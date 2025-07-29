package generator

import (
	"fmt"
	"log"
	"time"

	"github.com/Duckademic/schedule-generator/types"
)

type ScheduleGeneratorConfig struct {
	LessonsValue       int
	Start              time.Time
	End                time.Time
	WorkLessons        []int // ПОЧАТОК З НЕДІЛІ нд пн вт ср чт пт сб
	MaxStudentWorkload int   // максимальна кількість пар для студентів на день
}

type ScheduleGenerator struct {
	ScheduleGeneratorConfig
	BusyGrid            [][]bool
	teacherService      TeacherService
	studentGroupService StudentGroupService
	studyLoadService    StudyLoadService
	lessonService       LessonService
	disciplineService   DisciplineService
	boneWeek            int
}

func NewScheduleGenerator(cfg ScheduleGeneratorConfig) (*ScheduleGenerator, error) {
	if len(cfg.WorkLessons) != 7 {
		return nil, fmt.Errorf("length of WorkLessons %d instead of 7", len(cfg.WorkLessons))
	}
	if cfg.Start.After(cfg.End) {
		return nil, fmt.Errorf("start date comes after end")
	}

	scheduleGenerator := ScheduleGenerator{
		ScheduleGeneratorConfig: cfg,
	}

	for date := cfg.Start; !date.After(cfg.End); date = date.AddDate(0, 0, 1) {
		scheduleGenerator.BusyGrid = append(scheduleGenerator.BusyGrid, make([]bool, cfg.WorkLessons[date.Weekday()]))
	}

	ls, err := NewLessonService(cfg.LessonsValue)
	if err != nil {
		return nil, err
	}
	scheduleGenerator.lessonService = ls

	return &scheduleGenerator, nil
}

func (g *ScheduleGenerator) SetTeachers(teachers []types.Teacher) error {
	ts, err := NewTeacherService(teachers, g.BusyGrid)
	if err != nil {
		return err
	}

	g.teacherService = ts
	return nil
}

func (g *ScheduleGenerator) SetStudentGroups(studentGroups []types.StudentGroup) error {
	sgs, err := NewStudentGroupService(studentGroups, g.MaxStudentWorkload, g.BusyGrid)
	if err != nil {
		return err
	}

	g.studentGroupService = sgs
	return nil
}

func (g *ScheduleGenerator) SetDisciplines(disciplines []types.Discipline) error {
	ds, err := NewDisciplineService(disciplines)
	if err != nil {
		return err
	}

	g.disciplineService = ds
	return nil
}

func (g *ScheduleGenerator) SetStudyLoads(studyLoads []types.StudyLoad) error {
	err := g.CheckServices([]bool{true, true, true})
	if err != nil {
		return err
	}

	sls, err := NewStudyLoadService(studyLoads, g.teacherService, g.studentGroupService, g.disciplineService)
	if err != nil {
		return err
	}

	g.studyLoadService = sls
	return nil
}

func (g *ScheduleGenerator) CheckStudyLoadService() error {
	if g.studyLoadService == nil {
		return fmt.Errorf("study load not set")
	}
	return nil
}

// 0 - teacher, 1 - student group, 2 - discipline, 3 - study load
func (g *ScheduleGenerator) CheckServices(services []bool) error {
	checks := append(services, make([]bool, 4-len(services))...)

	if checks[0] && g.teacherService == nil {
		return fmt.Errorf("teachers not set")
	}

	if checks[1] && g.studentGroupService == nil {
		return fmt.Errorf("student groups not set")
	}

	if checks[2] && g.disciplineService == nil {
		return fmt.Errorf("discipline not set")
	}

	if checks[3] && g.studyLoadService == nil {
		return fmt.Errorf("study load not set")
	}

	return nil
}

func (g *ScheduleGenerator) GenerateSchedule() error {
	err := g.CheckServices([]bool{true, true, true, true})
	if err != nil {
		return err
	}

	err = g.generateBoneLectures()
	if err != nil {
		return err
	}

	g.buildLessonCarcass()

	err = g.addMissingLessons()
	if err != nil {
		return err
	}

	return nil
}

func (g *ScheduleGenerator) generateBoneLectures() error {
	for _, studyLoad := range g.studyLoadService.GetAll() {
		for _, dp := range studyLoad.Disciplines {
			for _, studentGroup := range dp.Groups {
				offset := 0
				success := false

				for !success {
					// отримуємо доступний лекційний день
					day := studentGroup.GetLectureDay(g.boneWeek*7 + offset)
					if day > g.boneWeek*7+7 {
						// якщо день був не на кістковому тижні, виникає виняток, який треба обробити якось
						return fmt.Errorf("group haven't enough slots for lectures")
					}

					// отримання вільного слота для групи та викладача
					lessonSlot := GetFirstFreeSlotForBoth(
						studentGroup.GetFreeSlots(day),
						studyLoad.Teacher.GetFreeSlots(day),
					)

					if lessonSlot != -1 {
						slot := LessonSlot{Day: day, Slot: lessonSlot}
						g.lessonService.CreateWithoutChecks(studyLoad.Teacher, studentGroup, dp.Discipline, slot, &LessonType{})
						success = true
					}
					offset = day - g.boneWeek*7 + 1
				}

			}
		}
	}

	return nil
}

func (g *ScheduleGenerator) buildLessonCarcass() {
	boneLessons := g.lessonService.GetWeekLessons(g.boneWeek)
	currentWeek := g.boneWeek + 1
	outOfGrid := false
	for !outOfGrid {
		for _, lesson := range boneLessons {
			newSlot := LessonSlot{
				Day:  lesson.Slot.Day + currentWeek*7,
				Slot: lesson.Slot.Slot,
			}

			err := g.lessonService.CreateWithChecks(
				lesson.Teacher,
				lesson.StudentGroup,
				lesson.Discipline,
				newSlot,
				lesson.Type,
			)
			if _, ok := err.(DayOutError); ok {
				outOfGrid = true
			}
		}
		currentWeek++
	}
}

func (g *ScheduleGenerator) addMissingLessons() error {
	for _, studyLoad := range g.studyLoadService.GetAll() {
		for _, disciplineLoad := range studyLoad.Disciplines {
			for _, group := range disciplineLoad.Groups {
				currentDay := g.boneWeek * 7
				outOfGrid := false
				for !disciplineLoad.Discipline.EnoughHours() && !outOfGrid {
					err := group.CheckDay(currentDay)
					if err != nil {
						outOfGrid = true
						//continue
						break
					}

					for i := range g.BusyGrid[currentDay] {
						slot := LessonSlot{
							Day:  currentDay,
							Slot: i,
						}
						g.lessonService.CreateWithChecks(
							studyLoad.Teacher,
							group,
							disciplineLoad.Discipline,
							slot,
							&LessonType{},
						)
					}
					currentDay++
				}

				if !disciplineLoad.Discipline.EnoughHours() {
					return fmt.Errorf("not enough space for %s discipline", disciplineLoad.Discipline.Name)
				}
			}
		}
	}

	return nil
}

func (g *ScheduleGenerator) CheckSchedule() error {
	err := g.CheckServices([]bool{true, true})
	if err != nil {
		return err
	}

	tw := g.teacherService.CountWindows()
	sgw := g.studentGroupService.CountWindows()
	log.Printf("вікна у викладачів: %d, вінка у студентів: %d", tw, sgw)
	hd := g.disciplineService.CountHourDeficit()
	lc := len(g.lessonService.GetAll())
	log.Printf("кількість занять: %d, недостача годин: %d", lc, hd)
	return nil
}

func (g *ScheduleGenerator) WriteSchedule() {
	// for _, l := range g.lessonService.GetAll() {
	// 	log.Printf("Generator викладач: %s, дисципліна: %s, група: %s, день/слот: %d/%d \n",
	// 		l.Teacher.UserName, l.Discipline.Name, l.StudentGroup.Name, l.Slot.Day, l.Slot.Slot,
	// 	)
	// }
	g.teacherService.WriteSchedule()
	g.studentGroupService.WriteSchedule()
}

func GetFirstFreeSlotForBoth(first, second []bool) int {
	min := min(len(first), len(second))
	for i := range min {
		if (first[i] == second[i]) && first[i] {
			return i
		}
	}
	return -1
}
