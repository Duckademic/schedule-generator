package generator

import (
	"fmt"
	"log"
	"slices"
	"time"

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
	teacherService      TeacherService
	studentGroupService StudentGroupService
	lessonService       LessonService
	disciplineService   DisciplineService
	lessonTypeService   LessonTypeService
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

			studentGroups := make([]*StudentGroup, len(disciplineLoad.GroupsID))
			for j, studentGroupID := range disciplineLoad.GroupsID {
				studentGroup := g.studentGroupService.Find(studentGroupID)
				if studentGroup == nil {
					return fmt.Errorf("student group %s not found", studentGroupID)
				}
				studentGroup.AddDayType(lessonType, disciplineLoad.Hours)

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

// 0 - teacher, 1 - student group, 2 - discipline, 3 - lesson type service
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

func (g *ScheduleGenerator) GenerateSchedule() error {
	if !g.studyLoadSet {
		return fmt.Errorf("study loads not set")
	}

	err := g.setDayTypes()
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
	studentGroups := g.studentGroupService.GetAll()

	type groupExtension struct {
		group         *StudentGroup
		dayPriorities []float32
		countOfSlots  []int
		freeDayCount  int
	}

	newGroupExtension := func(group *StudentGroup) *groupExtension {
		ge := groupExtension{
			group:         group,
			dayPriorities: group.GetWeekDaysPriority(),
			countOfSlots:  make([]int, 7),
		}

		for day, value := range ge.dayPriorities {
			if value > 1 {
				ge.freeDayCount++
			}
			ge.countOfSlots[day] = ge.group.CountSlotsAtDay(day)
		}

		return &ge
	}

	groupExtensions := make([]groupExtension, len(studentGroups))
	for i := range studentGroups {
		groupExtensions[i] = *newGroupExtension(&studentGroups[i])
	}
	slices.SortFunc(groupExtensions, func(a, b groupExtension) int {
		if a.freeDayCount == b.freeDayCount {
			return 0
		} else if a.freeDayCount > b.freeDayCount {
			return 1
		}
		return -1
	})

	lessonTypes := g.lessonTypeService.GetAll()
	for lessonTypeIndex := range lessonTypes {
		lType := &lessonTypes[lessonTypeIndex]
		dayOccupationsCount := make([]int, 7)

		for _, group := range groupExtensions {
			currentMaxHours := 0
			availableDays := []int{0, 1, 2, 3, 4, 5, 6}

			for currentMaxHours < group.group.GetMaxHours(lType) {
				min := 1000000000
				mIndex := -1
				for j := range availableDays {
					if min > dayOccupationsCount[j] && group.dayPriorities[j] > 1.0 {
						min = dayOccupationsCount[j]
						mIndex = j
					}
				}

				if mIndex == -1 {
					return fmt.Errorf("can't add a day of type %s to group %s", lType.Name, group.group.Name)
				}
				err := group.group.SetDayType(lType, mIndex)
				if err != nil {
					dayIndex := slices.Index(availableDays, mIndex)
					availableDays = append(availableDays[:dayIndex], availableDays[dayIndex+1:]...)
					continue
				}
				dayOccupationsCount[mIndex]++
				currentMaxHours += int(float64(group.countOfSlots[mIndex]*g.LessonsValue) * g.FillPercentage)
			}
		}
	}

	return nil
}

func (g *ScheduleGenerator) generateBoneLectures() error {
	teachers := g.teacherService.GetAll()

	for i := range teachers {
		teacher := &teachers[i]

		for _, teacherLoad := range teacher.Load {
			for _, studentGroup := range teacherLoad.Groups {
				offset := 0
				success := false

				for !success {
					// отримуємо доступний лекційний день
					day := studentGroup.GetNextDayOfType(teacherLoad.LessonType, g.boneWeek*7+offset)
					if day > g.boneWeek*7+7 || day < 0 {
						break
						// якщо день був не на кістковому тижні, виникає виняток, який треба обробити якось
						// return fmt.Errorf("group haven't enough slots for lectures")
					}

					// отримання вільного слота для групи та викладача
					lessonSlot := teacher.GetFreeSlot(studentGroup.GetFreeSlots(day), day)

					if lessonSlot != -1 {
						slot := LessonSlot{Day: day, Slot: lessonSlot}
						g.lessonService.AddWithoutChecks(teacher, studentGroup, teacherLoad.Discipline, slot, teacherLoad.LessonType)
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

			err := g.lessonService.AddWithChecks(
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
	teachers := g.teacherService.GetAll()

	for i := range teachers {
		teacher := &teachers[i]

		for _, teacherLoad := range teacher.Load {
			for _, group := range teacherLoad.Groups {
				currentDay := g.boneWeek * 7
				outOfGrid := false
				for !teacherLoad.Discipline.EnoughHours() && !outOfGrid {
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
						g.lessonService.AddWithChecks(
							teacher,
							group,
							teacherLoad.Discipline,
							slot,
							teacherLoad.LessonType,
						)
					}
					currentDay = group.GetNextDayOfType(teacherLoad.LessonType, currentDay+1)
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

	teacherWindows := g.teacherService.CountWindows()
	studentGroupWindows := g.studentGroupService.CountWindows()
	log.Printf("вікна у викладачів: %d, вінка у студентів: %d", teacherWindows, studentGroupWindows)

	disciplineHourDeficit := g.disciplineService.CountHourDeficit()
	lessonsCount := len(g.lessonService.GetAll())
	log.Printf("кількість занять: %d, недостача годин для дисциплін: %d", lessonsCount, disciplineHourDeficit)

	studentGroupHourDeficit := g.studentGroupService.CountHourDeficit()
	teachersHourDeficit := g.teacherService.CountHourDeficit()
	log.Printf("недостача годин для груп: %d, недостача годин для викладачів: %d", studentGroupHourDeficit, teachersHourDeficit)

	studentGroupLessonOverlapping := g.studentGroupService.CountLessonOverlapping()
	teacherLessonOverlapping := g.teacherService.CountLessonOverlapping()
	log.Printf("перекриття у студентів: %d, перекриття у викладачів: %d", studentGroupLessonOverlapping, teacherLessonOverlapping)

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
