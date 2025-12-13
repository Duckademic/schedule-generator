package generator

import (
	"fmt"
	"time"

	"github.com/Duckademic/schedule-generator/generator/components"
	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/generator/services"
	"github.com/Duckademic/schedule-generator/types"
)

type ScheduleGeneratorConfig struct {
	LessonsValue       int
	Start              time.Time
	End                time.Time
	WorkLessons        [][]float32 // ПОЧАТОК З НЕДІЛІ нд пн вт ср чт пт сб, зберігає коефіцієнти зручності
	MaxStudentWorkload int         // максимальна кількість пар для студентів на день
	FillPercentage     float64     // відсоток заповненості типом пар для визначення кількості днів
}

type ScheduleGenerator struct {
	ScheduleGeneratorConfig
	BusyGrid            [][]float32
	teacherService      services.TeacherService
	studentGroupService services.StudentGroupService
	lessonService       services.LessonService
	disciplineService   services.DisciplineService
	lessonTypeService   services.LessonTypeService
	errorService        components.ErrorService
	boneWeek            int
	studyLoadSet        bool
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

	ls, err := services.NewLessonService(cfg.LessonsValue)
	if err != nil {
		return nil, err
	}
	scheduleGenerator.lessonService = ls

	scheduleGenerator.errorService = components.NewErrorService()

	return &scheduleGenerator, nil
}

func (g *ScheduleGenerator) SetTeachers(teachers []types.Teacher) error {
	ts, err := services.NewTeacherService(teachers, g.BusyGrid)
	if err != nil {
		return err
	}

	g.teacherService = ts
	return nil
}

func (g *ScheduleGenerator) SetStudentGroups(studentGroups []types.StudentGroup) error {
	sgs, err := services.NewStudentGroupService(studentGroups, g.MaxStudentWorkload, g.BusyGrid)
	if err != nil {
		return err
	}

	g.studentGroupService = sgs
	return nil
}

func (g *ScheduleGenerator) SetDisciplines(disciplines []types.Discipline) error {
	ds, err := services.NewDisciplineService(disciplines)
	if err != nil {
		return err
	}

	g.disciplineService = ds
	return nil
}

func (g *ScheduleGenerator) SetLessonTypes(lTypes []types.LessonType) error {
	lts, err := services.NewLessonTypeService(lTypes)
	if err != nil {
		return err
	}

	g.lessonTypeService = lts
	g.boneWeek = g.lessonTypeService.GetWeekOffset()
	return nil
}

func (g *ScheduleGenerator) SetStudyLoads(studyLoads []types.StudyLoad) error {
	err := g.CheckServices([]bool{true, true, true, true})
	if err != nil {
		return err
	}

	for _, studyLoad := range studyLoads {
		teacher := g.teacherService.Find(studyLoad.TeacherID)
		if teacher == nil {
			return fmt.Errorf("teacher %s not found", studyLoad.TeacherID)
		}

		for _, disciplineLoad := range studyLoad.Disciplines {
			discipline := g.disciplineService.Find(disciplineLoad.DisciplineID)
			if discipline == nil {
				return fmt.Errorf("discipline %s not found", disciplineLoad.DisciplineID)
			}
			lessonType := g.lessonTypeService.Find(disciplineLoad.LessonTypeID)
			if lessonType == nil {
				return fmt.Errorf("lesson type %s not found", disciplineLoad.LessonTypeID)
			}

			studentGroups := make([]*entities.StudentGroup, len(disciplineLoad.GroupsID))
			for j, studentGroupID := range disciplineLoad.GroupsID {
				studentGroup := g.studentGroupService.Find(studentGroupID)
				if studentGroup == nil {
					return fmt.Errorf("student group %s not found", studentGroupID)
				}
				studentGroup.AddBindingToLessonType(lessonType, disciplineLoad.Hours, teacher)
				for week := range lessonType.Weeks {
					studentGroup.AddWeekToLessonType(lessonType, week)
				}

				studentGroups[j] = studentGroup
			}

			if err := discipline.AddLoad(teacher, disciplineLoad.Hours, studentGroups, lessonType); err != nil {
				return err
			}
			if err := teacher.AddLoad(discipline, lessonType, studentGroups, disciplineLoad.Hours); err != nil {
				return err
			}
		}
	}

	g.studyLoadSet = true
	return nil
}

// 0 - teacher, 1 - student group, 2 - discipline, 3 - lesson type service.
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

	if checks[3] && g.lessonTypeService == nil {
		return fmt.Errorf("lesson types not set")
	}

	return nil
}

// main function
func (g *ScheduleGenerator) GenerateSchedule() error {
	if !g.studyLoadSet {
		return fmt.Errorf("study loads not set")
	}

	components.NewDayBlocker(g.studentGroupService.GetAll(), g.errorService).SetDayTypes()

	components.NewBoneGenerator(g.errorService, g.teacherService.GetAll(), g.lessonService, g.boneWeek).GenerateBoneLessons()

	g.buildLessonCarcass()

	components.NewMissingLessonAdder(g.errorService, g.teacherService.GetAll(), g.lessonService).AddMissingLessons()

	// improver := components.NewImprover(g.lessonService)
	// improver.SubmitChanges() // CRUNCH - sets start slots for first lesson
	// result := true
	// currentFault := g.ScheduleFault()
	// for result {
	// 	if currentFault <= 0.000001 {
	// 		break
	// 	}
	// 	fault := g.ScheduleFault()
	// 	if fault < currentFault {
	// 		improver.SubmitChanges()
	// 	}
	// 	result = improver.ImproveToNext()
	// }

	if !g.errorService.IsClear() {
		return g.errorService
	}
	return nil
}

func (g *ScheduleGenerator) buildLessonCarcass() {
	boneLessons := g.lessonService.GetWeekLessons(g.boneWeek)
	currentWeek := 0
	outOfGrid := false
	for !outOfGrid && currentWeek <= len(g.BusyGrid)/7+1 {
		for _, lesson := range boneLessons {
			newSlot := entities.LessonSlot{
				Day:  lesson.Slot.Day%7 + currentWeek*7,
				Slot: lesson.Slot.Slot,
			}

			err := g.lessonService.AddLesson(
				lesson.Teacher,
				lesson.StudentGroup,
				lesson.Discipline,
				newSlot,
				lesson.Type,
			)
			if _, ok := err.(entities.DayOutError); ok {
				outOfGrid = true
			}
		}
		currentWeek++
	}
}

// Rates schedule fault. Returns ScheduleFault as a result.
// Returns an empty ScheduleFault if an not enough data.
func (g *ScheduleGenerator) ScheduleFault() (result components.ScheduleFault) {
	result = components.NewScheduleFault()

	err := g.CheckServices([]bool{true, true})
	if err != nil {
		return
	}

	result.AddParameter("teacher_windows", components.NewSimpleScheduleParameter(
		float64(g.teacherService.CountWindows()), 0.1,
	))
	result.AddParameter("student_group_windows", components.NewSimpleScheduleParameter(
		float64(g.studentGroupService.CountWindows()), 1000,
	))
	result.AddParameter("hours_deficit", components.NewSimpleScheduleParameter(
		float64(g.disciplineService.CountHourDeficit()), 10,
	))
	result.AddParameter("teacher_lesson_overlapping", components.NewSimpleScheduleParameter(
		float64(g.teacherService.CountLessonOverlapping()), 10,
	))
	result.AddParameter("student_group_lesson_overlapping", components.NewSimpleScheduleParameter(
		float64(g.studentGroupService.CountLessonOverlapping()), 10,
	))
	result.AddParameter("student_group_overtime_lessons", components.NewSimpleScheduleParameter(
		float64(g.studentGroupService.CountOvertimeLessons()), 10,
	))
	result.AddParameter("student_group_invalid_lessons_by_type", components.NewSimpleScheduleParameter(
		float64(g.studentGroupService.CountInvalidLessonsTypes()), 10,
	))

	return
}

func (g *ScheduleGenerator) WriteSchedule() {
	// for _, l := range g.lessonService.GetAll() {
	// 	log.Printf("Generator викладач: %s, дисципліна: %s, група: %s, день/слот: %d/%d \n",
	// 		l.Teacher.UserName, l.Discipline.Name, l.StudentGroup.Name, l.Slot.Day, l.Slot.Slot,
	// 	)
	// }
	tSchedule := make(map[*entities.Teacher]*entities.PersonalSchedule, len(g.teacherService.GetAll()))
	for i := range g.teacherService.GetAll() {
		t := g.teacherService.GetAll()[i]
		tSchedule[t] = &entities.PersonalSchedule{
			BusyGrid: &t.BusyGrid,
			Out:      "schedule/" + t.UserName + ".txt",
		}
	}

	sgSchedule := make(map[*entities.StudentGroup]*entities.PersonalSchedule, len(g.studentGroupService.GetAll()))
	for i := range g.studentGroupService.GetAll() {
		sg := g.studentGroupService.GetAll()[i]
		sgSchedule[sg] = &entities.PersonalSchedule{
			BusyGrid: &sg.BusyGrid,
			Out:      "schedule/" + sg.Name + ".txt",
		}
	}

	for _, l := range g.lessonService.GetAll() {
		tSchedule[l.Teacher].InsertLesson(l)
		sgSchedule[l.StudentGroup].InsertLesson(l)
	}

	for _, ps := range tSchedule {
		ps.WritePS(func(l *entities.Lesson) string {
			return fmt.Sprintf("дисципліна: %s, тип: %s, група: %s", l.Discipline.Name, l.Type.Name, l.StudentGroup.Name)
		})
	}
	for _, ps := range sgSchedule {
		ps.WritePS(func(l *entities.Lesson) string {
			return fmt.Sprintf("дисципліна: %s, тип: %s, викладач: %s", l.Discipline.Name, l.Type.Name, l.Teacher.UserName)
		})
	}
}
