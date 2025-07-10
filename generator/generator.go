package generator

import (
	"fmt"
	"log"
	"time"

	"github.com/Duckademic/schedule-generator/services"
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

func (g *ScheduleGenerator) GenerateShadule(
	studentGroupService *services.StudentGroupServise,
	teacherService *services.TeacherService,
	studyLoads []types.StudyLoad,
) (lessons []types.Lesson) {
	teacherService.SetBusyness(g.Business)
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
					lessonSlot := GetFirstPeretun(
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

// НЕТЕСТОВАНА + перейменувати
func GetFirstPeretun(first, second []bool) int {
	min := min(len(first), len(second))
	for i := range min {
		if (first[i] == second[i]) && first[i] {
			return i
		}
	}
	return -1
}
