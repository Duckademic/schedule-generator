package services

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type TeacherService interface {
	Find(uuid.UUID) *entities.Teacher
	GetAll() []*entities.Teacher
	CountWindows() int
	CountHourDeficit() int
	CountLessonOverlapping() int
}

type teacherService struct {
	teachers []*entities.Teacher
}

func NewTeacherService(teachers []types.Teacher, busyGrid [][]float32) (TeacherService, error) {
	ts := teacherService{teachers: make([]*entities.Teacher, 0, len(teachers))}

	for i := range teachers {
		teacher := &entities.Teacher{ID: teachers[i].ID, UserName: teachers[i].UserName, Priority: teachers[i].Priority}
		teacher.BusyGrid = *entities.NewBusyGrid(busyGrid)
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
				ts.teachers = append(ts.teachers[:j], append([]*entities.Teacher{teacher}, ts.teachers[j:]...)...)
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

func (ts *teacherService) GetAll() []*entities.Teacher {
	return ts.teachers
}

// return will be nil if not found
func (ts *teacherService) Find(id uuid.UUID) *entities.Teacher {
	for i := range ts.teachers {
		if ts.teachers[i].ID == id {
			return ts.teachers[i]
		}
	}

	return nil
}

// Returns sum of all teachers windows
// Time complexity O(n)
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

// Returns sum of all lesson overlap.
// Time complexity O(n^2)
func (ts *teacherService) CountLessonOverlapping() (count int) {
	for _, teacher := range ts.teachers {
		count += teacher.CountLessonOverlapping(teacher.Lessons)
	}

	return
}
