package controllers

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/services"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TeacherController interface {
	Controller[types.Teacher]
}

type teacherController struct {
	service  services.TeacherService
	validate *validator.Validate
}

func NewTeacherController(s services.TeacherService) TeacherController {
	return &teacherController{service: s, validate: validator.New()}
}

func (tc *teacherController) getTeahcerFromContext(ctx *gin.Context) (*types.Teacher, error) {
	var teacher types.Teacher
	err := ctx.ShouldBindBodyWithJSON(&teacher)
	if err != nil {
		return nil, err
	}

	// err = tc.validate.Struct(teacher)
	// if err != nil {
	// 	return nil, err
	// }

	return &teacher, nil
}

func (tc *teacherController) Create(ctx *gin.Context) (*types.Teacher, error) {
	teacher, err := tc.getTeahcerFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return tc.service.Create(*teacher)
}

func (tc *teacherController) Update(ctx *gin.Context) error {
	teacher, err := tc.getTeahcerFromContext(ctx)
	if err != nil {
		return err
	}

	return tc.service.Update(*teacher)
}

func (tc *teacherController) Delete(ctx *gin.Context) error {
	teacherId, ok := ctx.Params.Get("teacher_id")
	if !ok {
		return fmt.Errorf("missing teacher_id in URL parameters")
	}

	teachersUUID, err := uuid.Parse(teacherId)
	if err != nil {
		return fmt.Errorf("incorrect teacher id")
	}

	return tc.service.Delete(teachersUUID)
}

func (tc *teacherController) GetAll() []types.Teacher {
	return tc.service.GetAll()
}
