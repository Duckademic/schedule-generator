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

type generatorData struct {
	busyGrid            [][]float32
	teacherService      services.TeacherService
	studentGroupService services.StudentGroupService
	lessonService       services.LessonService
	disciplineService   services.DisciplineService
	lessonTypeService   services.LessonTypeService
	studyLoadService    services.StudyLoadService
}

// 0 - teacher, 1 - student group, 2 - discipline, 3 - lesson type service.
func (g *generatorData) CheckServices(services []bool) error {
	checks := append(services, make([]bool, 4-len(services))...)

	if checks[0] && g.teacherService == nil {
		return fmt.Errorf("teachers not set")
	}

	if checks[1] && g.studentGroupService == nil {
		return fmt.Errorf("student groups not set")
	}

	if checks[2] && g.disciplineService == nil {
		return fmt.Errorf("disciplines not set")
	}

	if checks[3] && g.lessonTypeService == nil {
		return fmt.Errorf("lesson types not set")
	}

	return nil
}

type ScheduleGenerator struct {
	ScheduleGeneratorConfig
	generatorData
	errorService components.ErrorService
	weekData     generatorData
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
		scheduleGenerator.busyGrid = append(scheduleGenerator.busyGrid, make([]float32, len(cfg.WorkLessons[date.Weekday()])))
		copy(scheduleGenerator.busyGrid[index], cfg.WorkLessons[date.Weekday()])
		index++
	}

	for i := range 7 {
		scheduleGenerator.weekData.busyGrid = append(scheduleGenerator.weekData.busyGrid,
			make([]float32, len(cfg.WorkLessons[i])))
		copy(scheduleGenerator.weekData.busyGrid[i], cfg.WorkLessons[i])
	}

	ls, err := services.NewLessonService(cfg.LessonsValue)
	if err != nil {
		return nil, err
	}

	weekLS, err := services.NewLessonService(cfg.LessonsValue)
	if err != nil {
		return nil, err
	}

	scheduleGenerator.lessonService = ls
	scheduleGenerator.weekData.lessonService = weekLS

	scheduleGenerator.errorService = components.NewErrorService()

	return &scheduleGenerator, nil
}

func (g *ScheduleGenerator) SetTeachers(teachers []types.Teacher) error {
	ts, err := services.NewTeacherService(teachers, g.busyGrid)
	if err != nil {
		return err
	}

	weekTS, err := services.NewTeacherService(teachers, g.weekData.busyGrid)
	if err != nil {
		return err
	}

	g.teacherService = ts
	g.weekData.teacherService = weekTS
	return nil
}

func (g *ScheduleGenerator) SetStudentGroups(studentGroups []types.StudentGroup) error {
	sgs, err := services.NewStudentGroupService(studentGroups, g.MaxStudentWorkload, g.busyGrid)
	if err != nil {
		return err
	}

	weekSGS, err := services.NewStudentGroupService(studentGroups, g.MaxStudentWorkload, g.weekData.busyGrid)
	if err != nil {
		return err
	}

	g.studentGroupService = sgs
	g.weekData.studentGroupService = weekSGS
	return nil
}

func (g *ScheduleGenerator) SetDisciplines(disciplines []types.Discipline) error {
	ds, err := services.NewDisciplineService(disciplines)
	if err != nil {
		return err
	}

	weekDS, err := services.NewDisciplineService(disciplines)
	if err != nil {
		return err
	}

	g.disciplineService = ds
	g.weekData.disciplineService = weekDS
	return nil
}

func (g *ScheduleGenerator) SetLessonTypes(lTypes []types.LessonType) error {
	lts, err := services.NewLessonTypeService(lTypes)
	if err != nil {
		return err
	}

	weekLTS, err := services.NewLessonTypeService(lTypes)
	if err != nil {
		return err
	}

	g.lessonTypeService = lts
	g.weekData.lessonTypeService = weekLTS
	return nil
}

func (g *ScheduleGenerator) SetStudyLoads(studyLoads []types.StudyLoad) error {
	if err := g.CheckServices([]bool{true, true, true, true}); err != nil {
		return err
	}
	if err := g.weekData.CheckServices([]bool{true, true, true, true}); err != nil {
		return err
	}

	sls, err := services.NewStudyLoadService(studyLoads, g.teacherService, g.studentGroupService,
		g.disciplineService, g.lessonTypeService)
	if err != nil {
		return err
	}

	weekSLS, err := services.NewStudyLoadService(studyLoads, g.weekData.teacherService, g.weekData.studentGroupService,
		g.weekData.disciplineService, g.weekData.lessonTypeService)
	if err != nil {
		return err
	}
	g.weekData.studentGroupService.UnbindWeeks()

	g.studyLoadService = sls
	g.weekData.studyLoadService = weekSLS
	return nil
}

// main function
func (g *ScheduleGenerator) GenerateSchedule() error {
	if g.studyLoadService == nil {
		return fmt.Errorf("study loads not set")
	}
	if g.weekData.studyLoadService == nil {
		return fmt.Errorf("study loads not set")
	}

	components.NewDayBlocker(g.weekData.studentGroupService.GetAll(), g.errorService).SetDayTypes()

	components.NewBoneGenerator(g.errorService, g.weekData.studyLoadService.GetAll(), g.weekData.lessonService).GenerateBoneLessons()
	g.buildLessonCarcass()

	// components.NewMissingLessonAdder(g.errorService, g.studyLoadService.GetAll(), g.lessonService).AddMissingLessons()

	// improver := components.NewImprover(g.lessonService)
	// improver.SubmitChanges() // CRUNCH - sets start slots for first lesson
	// result := true
	// currentFault := g.ScheduleFault()
	// for result {
	// 	if currentFault.Fault() <= 0.000001 {
	// 		break
	// 	}
	// 	fault := g.ScheduleFault()
	// 	if fault.Fault() < currentFault.Fault() {
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
	lessons := g.weekData.lessonService.GetAll()
	for _, lesson := range lessons {
		teacher := g.teacherService.Find(lesson.Teacher.ID)
		studentGroup := g.studentGroupService.Find(lesson.StudentGroup.ID)
		for weekday := range 7 {
			weekLT := lesson.StudentGroup.GetTypeOfDay(weekday)
			if weekLT != nil {
				lt := studentGroup.GetTypeOfDay(weekday)
				if lt == nil {
					lt := g.lessonTypeService.Find(weekLT.ID)
					err := studentGroup.BindWeekday(lt, weekday)
					if err != nil {
						g.errorService.AddError(components.NewUnexpectedError("can't bind the lesson type to the day",
							"ScheduleGenerator", "buildLessonCarcass", err))
					}
				}
			}
		}
		discipline := g.disciplineService.Find(lesson.Discipline.ID)
		lessonType := g.lessonTypeService.Find(lesson.Type.ID)

		currentWeek := 0
		outOfGrid := false
		for !outOfGrid {
			err := g.lessonService.AssignLesson(*entities.NewUnassignedLesson(lessonType, teacher, studentGroup, discipline),
				entities.NewLessonSlot(lesson.Day+currentWeek*7, lesson.Slot))
			if _, ok := err.(entities.DayOutError); ok {
				outOfGrid = true
			}
			currentWeek++
		}
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
	result.AddParameter("teacher_hours_deficit", components.NewSimpleScheduleParameter(
		float64(g.teacherService.CountHourDeficit()), 10,
	))
	result.AddParameter("student_group_hours_deficit", components.NewSimpleScheduleParameter(
		float64(g.studentGroupService.CountHourDeficit()), 10,
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
		float64(g.studentGroupService.CountInvalidLessonsByType()), 10,
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
