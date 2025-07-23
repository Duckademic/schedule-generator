package controllers

import (
	"github.com/Duckademic/schedule-generator/services"
	"github.com/Duckademic/schedule-generator/types"
)

type StudentGroupController interface {
	Controller[types.StudentGroup]
}

func NewStudentGroupController(s services.StudentGroupService) StudentGroupController {
	sgc := studentGroupController{
		basicController: basicController[types.StudentGroup]{
			service:       s,
			objectParamId: "student_group_id",
		},
		service: s,
	}

	return &sgc
}

type studentGroupController struct {
	basicController[types.StudentGroup]
	service services.StudentGroupService
}
