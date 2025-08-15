package services

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/repositories"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type TeacherService interface {
	Service[types.Teacher]
}

type teacherService struct {
	teachers []types.Teacher
}

func NewTeacherService(teachers []types.Teacher) TeacherService {
	ts := teacherService{teachers: teachers}
	return &ts
}

func (ts *teacherService) Create(teacher types.Teacher) (*types.Teacher, error) {
	t := ts.Find(teacher.ID)
	if t != nil {
		return nil, fmt.Errorf("teacher %s already exists", teacher.ID.String())
	}

	ts.teachers = append(ts.teachers, teacher)
	return &teacher, nil
}

func (ts *teacherService) Update(teacher types.Teacher) error {
	t := ts.Find(teacher.ID)
	if t == nil {
		return fmt.Errorf("teacher %s not found", teacher.ID.String())
	}

	t.UserName = teacher.UserName
	return nil
}

// return will be nil if not found
func (ts *teacherService) Find(id uuid.UUID) *types.Teacher {
	var teacher *types.Teacher
	for i := range ts.teachers {
		if ts.teachers[i].ID == id {
			teacher = &ts.teachers[i]
			break
		}
	}

	return teacher
}

func (ts *teacherService) Delete(teacherId uuid.UUID) error {
	for i, t := range ts.teachers {
		if t.ID == teacherId {
			ts.teachers = append(ts.teachers[:i], ts.teachers[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("teacher %s not found", teacherId)
}

func (ts *teacherService) GetAll() []types.Teacher {
	return ts.teachers
}

func NewGORMTeacherService(repo repositories.TeacherRepository) (TeacherService, error) {
	ts := gormTeacherService{
		gormSimpleService: gormSimpleService[types.Teacher]{
			repo: repo,
		},
	}

	return &ts, nil
}

type gormTeacherService struct {
	gormSimpleService[types.Teacher]
}
