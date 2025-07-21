package controllers

import (
	"github.com/Duckademic/schedule-generator/services"
	"github.com/Duckademic/schedule-generator/types"
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

type teacherController struct {
	basicController[types.Teacher]
	service services.TeacherService
}
