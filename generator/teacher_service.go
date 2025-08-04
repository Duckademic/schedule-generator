package generator

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type TeacherLoad struct {
	Discipline *Discipline
	Groups     []*StudentGroup
	LessonType *LessonType
}

type Teacher struct {
	BusyGrid
	LessonChecker
	ID       uuid.UUID
	UserName string
	Priority int
	Load     []TeacherLoad
}

func (t *Teacher) AddLoad(discipline *Discipline, lessonType *LessonType, groups []*StudentGroup, hours int) error {
	tl := TeacherLoad{
		Discipline: discipline,
		LessonType: lessonType,
		Groups:     groups,
	}

	t.RequiredHours += hours

	t.Load = append(t.Load, tl)
	return nil
}

func (t *Teacher) AddLesson(lesson *Lesson, ignoreCheck bool) error {
	err := t.CheckLesson(lesson)
	if err != nil && !ignoreCheck {
		return err
	}

	t.SetOneSlotBusyness(lesson.Slot, true)
	t.LessonChecker.AddLesson(lesson)

	return err
}

func (t *Teacher) CheckLesson(lesson *Lesson) error {
	if err := t.CheckSlot(lesson.Slot); err != nil {
		return err
	}
	if t.IsBusy(lesson.Slot) {
		return fmt.Errorf("teacher is busy")
	}
	if t.CountHourDeficit() <= 0 {
		return fmt.Errorf("enough hours")
	}

	return nil
}

type TeacherService interface {
	Find(uuid.UUID) *Teacher
	GetAll() []Teacher
	CountWindows() int
	CountHourDeficit() int
	CountLessonOverlapping() int
}

type teacherService struct {
	teachers []Teacher
}

func NewTeacherService(teachers []types.Teacher, busyGrid [][]float32) (TeacherService, error) {
	ts := teacherService{teachers: make([]Teacher, len(teachers))}

	for i := range teachers {
		teacher := Teacher{ID: teachers[i].ID, UserName: teachers[i].UserName, Priority: teachers[i].Priority}
		teacher.BusyGrid = *NewBusyGrid(busyGrid)

		success := false
		for j, lowerTeacher := range ts.teachers {
			if lowerTeacher.Priority <= teacher.Priority {
				ts.teachers = append(ts.teachers[:j], append([]Teacher{teacher}, ts.teachers[j:]...)...)
				success = true
				break
			}
		}
		if !success {
			ts.teachers = append(ts.teachers, teacher)
		}
	}

	return &ts, nil
}

func (ts *teacherService) GetAll() []Teacher {
	return ts.teachers
}

// return will be nil if not found
func (ts *teacherService) Find(id uuid.UUID) *Teacher {
	for i := range ts.teachers {
		if ts.teachers[i].ID == id {
			return &ts.teachers[i]
		}
	}

	return nil
}

func (ts *teacherService) CountWindows() (count int) {
	for _, t := range ts.teachers {
		count += t.CountWindows()
	}
	return
}

func (ts *teacherService) CountHourDeficit() (count int) {
	for _, teacher := range ts.teachers {
		count += teacher.CountHourDeficit()
	}

	return
}

func (ts *teacherService) CountLessonOverlapping() (count int) {
	for _, teacher := range ts.teachers {
		count += teacher.CountLessonOverlapping(teacher.Lessons)
	}

	return
}
