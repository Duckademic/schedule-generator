package services

import (
	"github.com/Duckademic/schedule-generator/types"
)

type TeacherService struct {
	teachers       []types.Teacher
	currentTeacher *types.Teacher
}

func NewTeacherService(teachers []types.Teacher) (*TeacherService, error) {
	ts := TeacherService{teachers: teachers}
	ts.currentTeacher = &ts.teachers[0]
	return &ts, nil
}

func (ts *TeacherService) GetAll() []types.Teacher {
	return ts.teachers
}

func (ts *TeacherService) SetBusyness(free [][]bool) {
	for i := range ts.teachers {
		ts.teachers[i].Business = make([][]bool, len(free))
		for j := range free {
			ts.teachers[i].Business[j] = make([]bool, len(free[j]))
			copy(ts.teachers[i].Business[j], free[j])
		}
	}
}

// НЕПРОТЕСТОВАНА ====================================================================
func (ts *TeacherService) SetOneSlotBusyness(teacherId string, day, slot int, isBusy bool) {
	teacher := ts.Find(teacherId)
	teacher.Business[day][slot] = isBusy
}

// НЕПРОТЕСТОВАНА ====================================================================
func (ts *TeacherService) GetFreeSlots(teacherId string, day int) (slots []bool) {
	teacher := ts.Find(teacherId)
	slots = make([]bool, len(teacher.Business[day]))

	// заглушка для викладача (вільні слоти - то всі незайняті)
	for i := range slots {
		slots[i] = !teacher.Business[day][i]
	}
	return
}

// НЕПРОТЕСТОВАНА ====================================================================
// return will be nil if not found
func (ts *TeacherService) Find(id string) *types.Teacher {
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
