package generator

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

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
		for _, day := range teachers[i].BusyDays {
			err := teacher.SetDayBusyness(make([]float32, len(busyGrid[day])), int(day))
			if err != nil {
				return nil, fmt.Errorf("teacher %s (%s) has invalid busy day %d (err: %s)",
					teacher.UserName, teacher.ID, day, err.Error(),
				)
			}
		}

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
