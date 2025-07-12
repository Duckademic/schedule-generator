package generator

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/types"
)

type TeacherService interface {
	ResetBusyness([][]bool)
	SetOneSlotBusyness(teacheId string, day, slot int, isBusy bool) error
	GetFreeSlots(teacherId string, day int) []bool
	Find(string) *types.Teacher
	GetAll() []types.Teacher
}

type teacherService struct {
	teachers       []types.Teacher
	currentTeacher *types.Teacher
}

func NewTeacherService(teachers []types.Teacher) (TeacherService, error) {
	ts := teacherService{teachers: teachers}
	ts.currentTeacher = &ts.teachers[0]
	return &ts, nil
}

func (ts *teacherService) GetAll() []types.Teacher {
	return ts.teachers
}

func (ts *teacherService) ResetBusyness(free [][]bool) {
	for i := range ts.teachers {
		ts.teachers[i].Business = make([][]bool, len(free))
		for j := range free {
			ts.teachers[i].Business[j] = make([]bool, len(free[j]))
			copy(ts.teachers[i].Business[j], free[j])
		}
	}
}

func (ts *teacherService) SetOneSlotBusyness(teacherId string, day, slot int, isBusy bool) error {
	if len(ts.teachers) == 0 {
		return fmt.Errorf("service hasn't teachers")
	}
	if len(ts.teachers[0].Business) <= day {
		return fmt.Errorf("day %d outside of the Business (%d)", day, len(ts.teachers[0].Business))
	}
	if len(ts.teachers[0].Business[day]) <= slot {
		return fmt.Errorf("teachers hasn't %d slot (max: %d)", slot, len(ts.teachers[0].Business[day]))
	}

	teacher := ts.Find(teacherId)
	if teacher == nil {
		return fmt.Errorf("teacher %s not found", teacherId)
	}

	teacher.Business[day][slot] = isBusy
	return nil
}

func (ts *teacherService) GetFreeSlots(teacherId string, day int) (slots []bool) {
	teacher := ts.Find(teacherId)
	slots = make([]bool, len(teacher.Business[day]))

	// заглушка для викладача (вільні слоти - то всі незайняті)
	for i := range slots {
		slots[i] = !teacher.Business[day][i]
	}
	return
}

// return will be nil if not found
func (ts *teacherService) Find(id string) *types.Teacher {
	var teacher *types.Teacher
	if ts.currentTeacher.UserName != id {
		for i := range ts.teachers {
			if ts.teachers[i].UserName == id {
				teacher = &ts.teachers[i]
				break
			}
		}
		ts.currentTeacher = teacher
	} else {
		teacher = ts.currentTeacher
	}

	return teacher
}
