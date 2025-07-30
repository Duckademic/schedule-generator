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
	WorkLessons        [][]float32 // ПОЧАТОК З НЕДІЛІ нд пн вт ср чт пт сб, зберігає коефіцієнти зручності
	MaxStudentWorkload int         // максимальна кількість пар для студентів на день
}

type ScheduleGenerator struct {
	ScheduleGeneratorConfig
	BusyGrid            [][]float32
	teacherService      TeacherService
	studentGroupService StudentGroupService
	studyLoadService    StudyLoadService
	lessonService       LessonService
	disciplineService   DisciplineService
	lessonTypeService   LessonTypeService
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

	index := 0
	for date := cfg.Start; !date.After(cfg.End); date = date.AddDate(0, 0, 1) {
		scheduleGenerator.BusyGrid = append(scheduleGenerator.BusyGrid, make([]float32, len(cfg.WorkLessons[date.Weekday()])))
		copy(scheduleGenerator.BusyGrid[index], cfg.WorkLessons[date.Weekday()])
		index++
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

func (g *ScheduleGenerator) SetLessonTypes(lTypes []types.LessonType) error {
	lts, err := NewLessonTypeService(lTypes)
	if err != nil {
		return err
	}

	g.lessonTypeService = lts
	return nil
}

func (g *ScheduleGenerator) SetStudyLoads(studyLoads []types.StudyLoad) error {
	err := g.CheckServices([]bool{true, true, true, false, true})
	if err != nil {
		return err
	}

	sls, err := NewStudyLoadService(studyLoads, g.teacherService, g.studentGroupService, g.disciplineService, g.lessonTypeService)
	if err != nil {
		return err
	}

	g.studyLoadService = sls
	return nil
}

// 0 - teacher, 1 - student group, 2 - discipline, 3 - study load, 4 - lesson type service
func (g *ScheduleGenerator) CheckServices(services []bool) error {
	checks := append(services, make([]bool, 5-len(services))...)

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

	if checks[4] && g.lessonTypeService == nil {
		return fmt.Errorf("lesson types not set")
	}

	return nil
}

func (g *ScheduleGenerator) GenerateSchedule() error {
	err := g.CheckServices([]bool{true, true, true, true})
	if err != nil {
		return err
	}

	err = g.setDayTypes()
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

func (g *ScheduleGenerator) setDayTypes() error {
	lessonTypes := g.lessonTypeService.GetAll()
	studentGroups := g.studentGroupService.GetAll()

	first := 0
	second := 1
	for i := range lessonTypes {
		for j := range studentGroups {
			studentGroups[j].SetDayType(&lessonTypes[i], first+1)
			studentGroups[j].SetDayType(&lessonTypes[i], second+1)

			first = second
			second = (second + 1) % 5
		}
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
					day := studentGroup.GetNextDayOfType(dp.LessonType, g.boneWeek*7+offset)
					if day > g.boneWeek*7+7 || day < 0 {
						// якщо день був не на кістковому тижні, виникає виняток, який треба обробити якось
						return fmt.Errorf("group haven't enough slots for lectures")
					}

					// отримання вільного слота для групи та викладача
					lessonSlot := studyLoad.Teacher.GetFreeSlot(studentGroup.GetFreeSlots(day), day)

					if lessonSlot != -1 {
						slot := LessonSlot{Day: day, Slot: lessonSlot}
						g.lessonService.CreateWithoutChecks(studyLoad.Teacher, studentGroup, dp.Discipline, slot, dp.LessonType)
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
							disciplineLoad.LessonType,
						)
					}
					currentDay = group.GetNextDayOfType(disciplineLoad.LessonType, currentDay+1)
				}

				// if !disciplineLoad.Discipline.EnoughHours() {
				// 	return fmt.Errorf("not enough space for %s discipline", disciplineLoad.Discipline.Name)
				// }
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
	tSchedule := make(map[*Teacher]*PersonalSchedule, len(g.teacherService.GetAll()))
	for i := range g.teacherService.GetAll() {
		t := &g.teacherService.GetAll()[i]
		tSchedule[t] = &PersonalSchedule{
			busyGrid: &t.BusyGrid,
			out:      "schedule/" + t.UserName + ".txt",
		}
	}

	sgSchedule := make(map[*StudentGroup]*PersonalSchedule, len(g.studentGroupService.GetAll()))
	for i := range g.studentGroupService.GetAll() {
		sg := &g.studentGroupService.GetAll()[i]
		sgSchedule[sg] = &PersonalSchedule{
			busyGrid: &sg.BusyGrid,
			out:      "schedule/" + sg.Name + ".txt",
		}
	}

	for _, l := range g.lessonService.GetAll() {
		tSchedule[l.Teacher].InsertLesson(&l)
		sgSchedule[l.StudentGroup].InsertLesson(&l)
	}

	for _, ps := range tSchedule {
		ps.WritePS(func(l *Lesson) string {
			return fmt.Sprintf("дисципліна: %s, тип: %s, група: %s", l.Discipline.Name, l.Type.Name, l.StudentGroup.Name)
		})
	}
	for _, ps := range sgSchedule {
		ps.WritePS(func(l *Lesson) string {
			return fmt.Sprintf("дисципліна: %s, тип: %s, викладач: %s", l.Discipline.Name, l.Type.Name, l.Teacher.UserName)
		})
	}
}
