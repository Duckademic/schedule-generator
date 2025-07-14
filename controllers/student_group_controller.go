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
	service  services.StudentGroupServise
	validate *validator.Validate
}

func NewStudentGroupController(service services.StudentGroupServise) StudentGroupController {
	sgc := studentGroupController{
		service:  service,
		validate: validator.New(),
	}

	return &sgc
}

func (sgc *studentGroupController) Create(ctx *gin.Context) (*types.StudentGroup, error) {
	studentGroup, err := sgc.getStudentGroupFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return sgc.service.Create(*studentGroup)
}

func (sgc *studentGroupController) Update(ctx *gin.Context) error {
	studentGroup, err := sgc.getStudentGroupFromContext(ctx)
	if err != nil {
		return err
	}

	return sgc.service.Update(*studentGroup)
}

func (sgc *studentGroupController) Delete(ctx *gin.Context) error {
	studentGroupID, ok := ctx.Params.Get("student_group_id")
	if !ok {
		return fmt.Errorf("missing student_group_id in URL parameters")
	}

	sgUUID, err := uuid.Parse(studentGroupID)
	if err != nil {
		return fmt.Errorf("incorrect student group id")
	}

	return sgc.service.Delete(sgUUID)
}

func (sgc *studentGroupController) GetAll() []types.StudentGroup {
	return sgc.service.GetAll()
}

func (sgc *studentGroupController) getStudentGroupFromContext(ctx *gin.Context) (*types.StudentGroup, error) {
	var studentGroup types.StudentGroup
	err := ctx.ShouldBindBodyWithJSON(&studentGroup)
	if err != nil {
		return nil, err
	}

	return &studentGroup, nil
}
