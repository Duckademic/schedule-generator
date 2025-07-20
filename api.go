package main

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/controllers"
	"github.com/Duckademic/schedule-generator/generator"
	"github.com/Duckademic/schedule-generator/services"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/gin-gonic/gin"
)

type JSONAPIServer struct {
	listenAddr             string
	generator              generator.ScheduleGenerator
	teacherController      controllers.TeacherController
	studentGroupController controllers.StudentGroupController
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
	api.studentGroupController = controllers.NewStudentGroupController(services.NewStudentGroupService([]types.StudentGroup{}))

	return &api, nil
}

func (s *JSONAPIServer) Run() error {
	server := gin.Default()

	// server.POST("/generator/reset/", func(ctx *gin.Context) {

	// })
	server.GET("/teacher/", s.teacherController.GetAll)
	server.POST("/teacher/", s.teacherController.Create)
	server.PUT("/teacher/:teacher_id/", s.teacherController.Update)
	server.DELETE("/teacher/:teacher_id/", s.teacherController.Delete)

	server.GET("/student_group/", s.studentGroupController.GetAll)
	server.POST("/student_group/", s.studentGroupController.Create)
	server.PUT("/student_group/:student_group_id/", s.studentGroupController.Update)
	server.DELETE("/student_group/:student_group_id/", s.studentGroupController.Delete)

	err := server.Run(s.listenAddr)
	return err
}
