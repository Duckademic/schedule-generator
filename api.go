package main

import (
	"fmt"
	"net/http"

	"github.com/Duckademic/schedule-generator/controllers"
	"github.com/Duckademic/schedule-generator/generator"
	"github.com/Duckademic/schedule-generator/services"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/gin-gonic/gin"
)

type JSONAPIServer struct {
	listenAddr        string
	generator         generator.ScheduleGenerator
	teacherController controllers.Controller[types.Teacher]
}

func NewJSONAPIServer(listenAddr string, cfg generator.ScheduleGeneratorConfig) (*JSONAPIServer, error) {
	gen, err := generator.NewScheduleGenerator(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't create generator: %s", err.Error())
	}

	api := JSONAPIServer{
		listenAddr: listenAddr,
		generator:  *gen,
	}

	api.teacherController = controllers.NewTeacherController(services.NewTeacherService([]types.Teacher{}))

	return &api, nil
}

func (s *JSONAPIServer) Run() error {
	server := gin.Default()

	// server.POST("/generator/reset/", func(ctx *gin.Context) {

	// })
	server.GET("/teacher/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, s.teacherController.GetAll())
	})

	server.POST("/teacher/", func(ctx *gin.Context) {
		teacher, err := s.teacherController.Create(ctx)
		if err != nil {
			s.responseWithError(ctx, http.StatusBadRequest, err)
			return
		}

		ctx.JSON(http.StatusCreated, teacher)
	})

	server.PUT("/teacher/:teacher_id/", func(ctx *gin.Context) {
		err := s.teacherController.Update(ctx)
		if err != nil {
			s.responseWithError(ctx, http.StatusBadRequest, err)
			return
		}

		ctx.Status(http.StatusNoContent)
	})

	server.DELETE("/teacher/:teacher_id/", func(ctx *gin.Context) {
		err := s.teacherController.Delete(ctx)
		if err != nil {
			s.responseWithError(ctx, http.StatusBadRequest, err)
			return
		}

		ctx.Status(http.StatusNoContent)
	})

	err := server.Run(s.listenAddr)
	return err
}

func (s *JSONAPIServer) responseWithError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, gin.H{"error": err.Error()})
}
