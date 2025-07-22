package generator

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type Teacher struct {
	ID       uuid.UUID
	UserName string
	BusyGrid [][]bool
}

type TeacherService interface {
	SetOneSlotBusyness(*Teacher, LessonSlot, bool) error
	GetFreeSlots(teacher *Teacher, day int) []bool
	Find(uuid.UUID) *Teacher
	GetAll() []Teacher
}

type teacherService struct {
	teachers []Teacher
}

func NewTeacherService(teachers []types.Teacher, busyGrid [][]bool) (TeacherService, error) {
	ts := teacherService{teachers: make([]Teacher, len(teachers))}

	for i := range teachers {
		ts.teachers[i] = Teacher{ID: teachers[i].ID, UserName: teachers[i].UserName}
		ts.teachers[i].BusyGrid = make([][]bool, len(busyGrid))
		for j := range busyGrid {
			ts.teachers[i].BusyGrid[j] = make([]bool, len(busyGrid[j]))
			copy(ts.teachers[i].BusyGrid[j], busyGrid[j])
		}
	}

	return &ts, nil
}

func (ts *teacherService) GetAll() []Teacher {
	return ts.teachers
}

func (ts *teacherService) SetOneSlotBusyness(teacher *Teacher, slot LessonSlot, isBusy bool) error {
	if len(ts.teachers) == 0 {
		return fmt.Errorf("service hasn't teachers")
	}
	if len(ts.teachers[0].BusyGrid) <= slot.Day {
		return fmt.Errorf("day %d outside of the Business (%d)", slot.Day, len(ts.teachers[0].BusyGrid))
	}
	if len(ts.teachers[0].BusyGrid[slot.Day]) <= slot.Slot {
		return fmt.Errorf("teachers hasn't %d slot (max: %d)", slot, len(ts.teachers[0].BusyGrid[slot.Day]))
	}
	if teacher == nil {
		return fmt.Errorf("teacher is nil")
	}

	teacher.BusyGrid[slot.Day][slot.Slot] = isBusy
	return nil
}

func (ts *teacherService) GetFreeSlots(teacher *Teacher, day int) (slots []bool) {
	if teacher == nil {
		return
	}

	slots = make([]bool, len(teacher.BusyGrid[day]))

	// заглушка для викладача (вільні слоти - то всі незайняті)
	for i := range slots {
		slots[i] = !teacher.BusyGrid[day][i]
	}
	return
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
