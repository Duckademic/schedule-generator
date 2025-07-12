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
	Business [][]bool
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
		scheduleGenerator.Business = append(scheduleGenerator.Business, make([]bool, cfg.WorkLessons[date.Weekday()]))
	}

	return &scheduleGenerator, nil
}

func (g *ScheduleGenerator) GenerateShedule(studyLoads []types.StudyLoad) error {
	studGroups := map[string]types.StudentGroup{}
	teachers := map[string]types.Teacher{}

	for _, sl := range studyLoads {
		teachers[sl.Teacher.UserName] = sl.Teacher
		for _, dl := range sl.Disciplines {
			for _, group := range dl.Groups {
				studGroups[group.Name] = group
			}
		}
	}

	studGroupsArr := make([]types.StudentGroup, len(studGroups))
	counter := 0
	for _, value := range studGroups {
		studGroupsArr[counter] = value
		counter++
	}
	studentGroupService, err := NewStudentGroupService(studGroupsArr, g.MaxStudentWorkload)
	if err != nil {
		return err
	}
	studentGroupService.SetBusyness(g.Business)

	teachersArr := make([]types.Teacher, len(teachers))
	counter = 0
	for _, value := range teachers {
		teachersArr[counter] = value
		counter++
	}
	teacherService, err := NewTeacherService(teachersArr)
	if err != nil {
		return err
	}
	teacherService.ResetBusyness(g.Business)

	log.Printf("%d %d \n", len(studGroupsArr), len(teachersArr))

	g.startGenerateShedule(studentGroupService, teacherService, studyLoads)

	return nil
}

func (g *ScheduleGenerator) startGenerateShedule(
	studentGroupService StudentGroupServise,
	teacherService TeacherService,
	studyLoads []types.StudyLoad,
) (lessons []types.Lesson) {
	teacherService.ResetBusyness(g.Business)
	studentGroupService.SetBusyness(g.Business)

	// номер кісткового тижня
	mainWeak := 0

	for _, studyLoad := range studyLoads {
		for _, dp := range studyLoad.Disciplines {
			for _, group := range dp.Groups {
				offset := 0
				success := false

				for !success {
					// отримуємо доступний лекційний день
					day := studentGroupService.GetLectureDay(group.Name, mainWeak*7+offset)
					if day > mainWeak*7+7 {
						// якщо день був не на кістковому тижні, виникає виняток, який треба обробити якось
						panic("group havn't enought slots for lectures")
					}

					// отримання вільного слота для групи та викладача
					lessonSlot := GetFirstFreeSlotForBoth(
						studentGroupService.GetFreeSlots(group.Name, day),
						teacherService.GetFreeSlots(studyLoad.Teacher.UserName, day),
					)

					if lessonSlot != -1 {
						// встановлення лекції
						// ЗАГЛУШКА
						log.Printf("викладач: %s, група: %s, день/слот: %d/%d \n", studyLoad.Teacher.UserName, group.Name, day, lessonSlot)
						// ========
						studentGroupService.SetOneSlotBusyness(group.Name, day, lessonSlot, true)
						teacherService.SetOneSlotBusyness(studyLoad.Teacher.UserName, day, lessonSlot, true)
						success = true
					}
					offset = day - mainWeak*7 + 1
				}

			}
		}
	}

	tw, sgw := types.CheckWindows(teacherService.GetAll(), studentGroupService.GetAll())
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
