package controllers

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/repositories"
	"github.com/Duckademic/schedule-generator/services"
	"github.com/Duckademic/schedule-generator/types"
	"gorm.io/gorm"
)

type TeacherController interface {
	Controller[types.Teacher]
}

func NewTeacherController(s services.TeacherService) TeacherController {
	tc := teacherController{
		basicController: basicController[types.Teacher]{
			service:       s,
			objectParamId: "teacher_id",
		},
		service: s,
	}

	return &tc
}

func NewDefaultTeacherController(db *gorm.DB) (TeacherController, error) {
	repo, err := repositories.NewTeacherRepository(db)
	if err != nil {
		return nil, fmt.Errorf("cannot crate teacher repository: %s", err)
	}

	s, err := services.NewGORMTeacherService(repo)
	if err != nil {
		return nil, fmt.Errorf("cannot create teacher service: %s", err)
	}

	return NewTeacherController(s), nil
}

type teacherController struct {
	basicController[types.Teacher]
	service services.TeacherService
}
