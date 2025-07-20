package controllers

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/services"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type StudentGroupController interface {
	Controller[types.StudentGroup]
}

type studentGroupController struct {
	BasicController[types.StudentGroup]
}

func NewStudentGroupController(service services.StudentGroupServise) StudentGroupController {
	sgc := studentGroupController{
		BasicController: BasicController[types.StudentGroup]{
			serviceController: &studentGroupServiceController{
				service:  service,
				validate: validator.New(),
			},
		},
	}

	return &sgc
}

type StudentGroupServiceController interface {
	ServiceController[types.StudentGroup]
}

type studentGroupServiceController struct {
	service  services.StudentGroupServise
	validate *validator.Validate
}

func (sgsc *studentGroupServiceController) Create(ctx *gin.Context) (*types.StudentGroup, error) {
	studentGroup, err := sgsc.getStudentGroupFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return sgsc.service.Create(*studentGroup)
}

func (sgsc *studentGroupServiceController) Update(ctx *gin.Context) error {
	studentGroup, err := sgsc.getStudentGroupFromContext(ctx)
	if err != nil {
		return err
	}

	return sgsc.service.Update(*studentGroup)
}

func (sgsc *studentGroupServiceController) Delete(ctx *gin.Context) error {
	studentGroupID, ok := ctx.Params.Get("student_group_id")
	if !ok {
		return fmt.Errorf("missing student_group_id in URL parameters")
	}

	sgUUID, err := uuid.Parse(studentGroupID)
	if err != nil {
		return fmt.Errorf("incorrect student group id")
	}

	return sgsc.service.Delete(sgUUID)
}

func (sgsc *studentGroupServiceController) GetAll() []types.StudentGroup {
	return sgsc.service.GetAll()
}

func (sgsc *studentGroupServiceController) getStudentGroupFromContext(ctx *gin.Context) (*types.StudentGroup, error) {
	var studentGroup types.StudentGroup
	err := ctx.ShouldBindBodyWithJSON(&studentGroup)
	if err != nil {
		return nil, err
	}

	return &studentGroup, nil
}
