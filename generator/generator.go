package generator

import (
	"fmt"
	"log"
	"time"

	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

// ==============================================================

type Discipline struct {
	ID   uuid.UUID
	Name string
	// Lessons map[string]int // тип - кількість годин
}

func CheckWindows(teachers []Teacher, groups []StudentGroup) (teacherW, groupW int) {
	for i := range len(teachers[0].BusyGrid) {
		for _, t := range teachers {
			lastBusy := -1
			for j, isBusy := range t.BusyGrid[i] {
				if isBusy {
					if lastBusy != -1 && (j-lastBusy) > 1 {
						teacherW += j - lastBusy - 1
					}
					lastBusy = j
				}
			}
		}

		for _, t := range groups {
			lastBusy := -1
			for j, isBusy := range t.BusyGrid[i] {
				if isBusy {
					if lastBusy != -1 && (j-lastBusy) > 1 {
						groupW += j - lastBusy - 1
					}
					lastBusy = j
				}
			}
		}
	}
	return
}

// ==================================================================================

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
	if len(g.BusyGrid) == 0 {
		return fmt.Errorf("config not initialized")
	}

	ts, err := NewTeacherService(teachers, g.BusyGrid)
	if err != nil {
		return err
	}

	g.teacherService = ts
	return nil
}

func (g *ScheduleGenerator) SetStudentGroups(studentGroups []types.StudentGroup) error {
	if len(g.BusyGrid) == 0 {
		return fmt.Errorf("config not initialized")
	}

	sgs, err := NewStudentGroupService(studentGroups, g.MaxStudentWorkload, g.BusyGrid)
	if err != nil {
		return err
	}

	g.studentGroupService = sgs
	return nil
}

func (g *ScheduleGenerator) SetStudyLoads(studyLoads []types.StudyLoad) error {
	if g.teacherService == nil {
		return fmt.Errorf("teachers not set")
	}
	if g.studentGroupService == nil {
		return fmt.Errorf("student groups not set")
	}

	sls, err := NewStudyLoadService(studyLoads, g.teacherService, g.studentGroupService)
	if err != nil {
		return err
	}

	g.studyLoadService = sls
	return nil
}

func (g *ScheduleGenerator) GenerateSchedule() error {
	if g.teacherService == nil {
		return fmt.Errorf("teachers not set")
	}
	if g.studentGroupService == nil {
		return fmt.Errorf("student groups not set")
	}
	if g.studyLoadService == nil {
		return fmt.Errorf("study load service not set")
	}

	// номер кісткового тижня
	mainWeak := 0

	for _, studyLoad := range g.studyLoadService.GetAll() {
		for _, dp := range studyLoad.Disciplines {
			for _, studentGroup := range dp.Groups {
				offset := 0
				success := false

				for !success {
					// отримуємо доступний лекційний день
					day := g.studentGroupService.GetLectureDay(studentGroup, mainWeak*7+offset)
					if day > mainWeak*7+7 {
						// якщо день був не на кістковому тижні, виникає виняток, який треба обробити якось
						panic("group haven't enough slots for lectures")
					}

					// отримання вільного слота для групи та викладача
					lessonSlot := GetFirstFreeSlotForBoth(
						g.studentGroupService.GetFreeSlots(studentGroup, day),
						g.teacherService.GetFreeSlots(studyLoad.Teacher, day),
					)

					if lessonSlot != -1 {
						slot := LessonSlot{Day: day, Slot: lessonSlot}
						g.lessonService.CreateWithoutChecks(studyLoad.Teacher, studentGroup, dp.Discipline, slot, LessonType{})
						g.studentGroupService.SetOneSlotBusyness(studentGroup, slot, true)
						g.teacherService.SetOneSlotBusyness(studyLoad.Teacher, slot, true)
						success = true
					}
					offset = day - mainWeak*7 + 1
				}

			}
		}
	}

	return nil
}

func (g *ScheduleGenerator) CheckSchedule() error {
	if g.teacherService == nil {
		return fmt.Errorf("teachers not set")
	}
	if g.studentGroupService == nil {
		return fmt.Errorf("student groups not set")
	}

	for _, l := range g.lessonService.GetAll() {
		log.Printf("викладач: %s, група: %s, день/слот: %d/%d \n",
			l.Teacher.UserName, l.StudentGroup.Name, l.Slot.Day, l.Slot.Slot,
		)
	}
	tw, sgw := CheckWindows(g.teacherService.GetAll(), g.studentGroupService.GetAll())
	log.Printf("вікна у викладачів: %d, вінка у студентів: %d", tw, sgw)
	return nil
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
