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

func NewTeacherController(s services.TeacherService) TeacherController {
	return &teacherController{
		BasicController: BasicController[types.Teacher]{
			serviceController: &teacherServiceController{
				service:  s,
				validate: validator.New(),
			},
		},
	}
}

type teacherController struct {
	BasicController[types.Teacher]
}

type TeacherServiceController interface {
	ServiceController[types.Teacher]
}

type teacherServiceController struct {
	service  services.TeacherService
	validate *validator.Validate
}

func (tsc *teacherServiceController) getTeahcerFromContext(ctx *gin.Context) (*types.Teacher, error) {
	var teacher types.Teacher
	err := ctx.ShouldBindBodyWithJSON(&teacher)
	if err != nil {
		return nil, err
	}

	// err = tsc.validate.Struct(teacher)
	// if err != nil {
	// 	return nil, err
	// }

	return &teacher, nil
}

func (tsc *teacherServiceController) Create(ctx *gin.Context) (*types.Teacher, error) {
	teacher, err := tsc.getTeahcerFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return tsc.service.Create(*teacher)
}

func (tsc *teacherServiceController) Update(ctx *gin.Context) error {
	teacher, err := tsc.getTeahcerFromContext(ctx)
	if err != nil {
		return err
	}

	return tsc.service.Update(*teacher)
}

func (tsc *teacherServiceController) Delete(ctx *gin.Context) error {
	teacherId, ok := ctx.Params.Get("teacher_id")
	if !ok {
		return fmt.Errorf("missing teacher_id in URL parameters")
	}

	teachersUUID, err := uuid.Parse(teacherId)
	if err != nil {
		return fmt.Errorf("incorrect teacher id")
	}

	return tsc.service.Delete(teachersUUID)
}

func (tsc *teacherServiceController) GetAll() []types.Teacher {
	return tsc.service.GetAll()
}
