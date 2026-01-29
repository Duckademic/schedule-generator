package services

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

// TeacherService aggregates and manages teachers that the generator works with.
type TeacherService interface {
	Find(uuid.UUID) *entities.Teacher // Returns a pointer to the teacher with the given ID.
	GetAll() []*entities.Teacher      // Returns a slice with all teachers as pointers.
	CountWindows() int                // Returns the sum of windows (gaps between busy slots).
	CountHourDeficit() int            // Returns the number of missing study hours.
	CountLessonOverlapping() int      // Returns the count of overlapping lessons.
}

// NewTeacherService creates a new TeacherService basic instance.
//
// It requires an array of database teachers (t) and a busy grid for them (bg).
//
// Returns an error if any teacher is an invalid model.
func NewTeacherService(t []types.Teacher, bg [][]float32) (TeacherService, error) {
	ts := teacherService{teachers: make([]*entities.Teacher, 0, len(t))}

	for i := range t {
		teacher := entities.NewDefaultTeacher(t[i].ID, t[i].UserName, t[i].Priority, entities.NewBusyGrid(bg))
		for _, day := range t[i].BusyDays {
			err := teacher.BlockWeekDay(int(day))
			if err != nil {
				return nil, fmt.Errorf("teacher %s (%s) has invalid busy day %d (err: %s)",
					teacher.UserName, teacher.ID, day, err.Error(),
				)
			}
		}

		// sort in priority order (not necessary now)
		// success := false
		// for j, lowerTeacher := range ts.teachers {
		// 	if lowerTeacher.Priority <= teacher.Priority {
		// 		ts.teachers = append(ts.teachers[:j], append([]*entities.Teacher{teacher}, ts.teachers[j:]...)...)
		// 		success = true
		// 		break
		// 	}
		// }
		// if !success {
		ts.teachers = append(ts.teachers, teacher)
		// }
	}

	return &ts, nil
}

// teacherService is the basic implementation of the TeacherService interface.
type teacherService struct {
	teachers []*entities.Teacher
}

func (ts *teacherService) GetAll() []*entities.Teacher {
	return ts.teachers
}
func (ts *teacherService) Find(id uuid.UUID) *entities.Teacher {
	for i := range ts.teachers {
		if ts.teachers[i].ID == id {
			return ts.teachers[i]
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
		count += teacher.CountLessonOverlapping(teacher.GetAssignedLessons())
	}

	return
}
