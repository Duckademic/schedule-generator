package generator

import (
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
	ID       uuid.UUID
	UserName string
	Load     []TeacherLoad
}

func (t *Teacher) AddLoad(tl *TeacherLoad) error {
	t.Load = append(t.Load, *tl)
	return nil
}

type TeacherService interface {
	Find(uuid.UUID) *Teacher
	GetAll() []Teacher
	CountWindows() int
}

type teacherService struct {
	teachers []Teacher
}

func NewTeacherService(teachers []types.Teacher, busyGrid [][]float32) (TeacherService, error) {
	ts := teacherService{teachers: make([]Teacher, len(teachers))}

	for i := range teachers {
		ts.teachers[i] = Teacher{ID: teachers[i].ID, UserName: teachers[i].UserName}
		ts.teachers[i].BusyGrid = *NewBusyGrid(busyGrid)
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
