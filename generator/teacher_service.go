package generator

import (
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type Teacher struct {
	ID       uuid.UUID
	UserName string
	BusyGrid
	PersonalSchedule
}

type TeacherService interface {
	Find(uuid.UUID) *Teacher
	GetAll() []Teacher
	CountWindows() int
	WriteSchedule()
}

type teacherService struct {
	teachers []Teacher
}

func NewTeacherService(teachers []types.Teacher, busyGrid [][]bool) (TeacherService, error) {
	ts := teacherService{teachers: make([]Teacher, len(teachers))}

	for i := range teachers {
		ts.teachers[i] = Teacher{ID: teachers[i].ID, UserName: teachers[i].UserName}
		ts.teachers[i].BusyGrid = *NewBusyGrid(busyGrid)
		ts.teachers[i].PersonalSchedule = PersonalSchedule{
			busyGrid: &ts.teachers[i].BusyGrid,
			out:      "schedule/" + ts.teachers[i].UserName + ".txt",
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

func (ts *teacherService) WriteSchedule() {
	for _, teacher := range ts.teachers {
		err := teacher.WritePS(StringForTeacher)
		if err != nil {
			panic(err)
		}
	}
}
