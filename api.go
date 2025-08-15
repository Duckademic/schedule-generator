package main

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/controllers"
	"github.com/Duckademic/schedule-generator/generator"
	"github.com/Duckademic/schedule-generator/services"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JSONAPIServer struct {
	listenAddr             string
	generator              generator.ScheduleGenerator
	teacherController      controllers.TeacherController
	studentGroupController controllers.StudentGroupController
	lessonController       controllers.LessonController
}

func NewJSONAPIServer(listenAddr string, cfg generator.ScheduleGeneratorConfig, db *gorm.DB) (*JSONAPIServer, error) {
	gen, err := generator.NewScheduleGenerator(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't create generator: %s", err.Error())
	}

	api := JSONAPIServer{
		listenAddr: listenAddr,
		generator:  *gen,
	}

	api.teacherController, err = controllers.NewDefaultTeacherController(db)
	if err != nil {
		return nil, fmt.Errorf("cannot create teacher controller: %s", err)
	}

	api.studentGroupController = controllers.NewStudentGroupController(services.NewStudentGroupService([]types.StudentGroup{}))
	api.lessonController = controllers.NewLessonController(services.NewLessonService([]types.Lesson{}))

	return &api, nil
}

func (s *JSONAPIServer) Run() error {
	server := gin.Default()

	// server.POST("/generator/reset/", func(ctx *gin.Context) {

	// })

	teacherRouts := server.Group("/teacher")
	teacherRouts.GET("/", s.teacherController.GetAll)
	teacherRouts.POST("/", s.teacherController.Create)
	teacherRouts.PUT("/:teacher_id/", s.teacherController.Update)
	teacherRouts.DELETE("/:teacher_id/", s.teacherController.Delete)

	studentGroupRouts := server.Group("/student_group")
	studentGroupRouts.GET("/", s.studentGroupController.GetAll)
	studentGroupRouts.POST("/", s.studentGroupController.Create)
	studentGroupRouts.PUT("/:student_group_id/", s.studentGroupController.Update)
	studentGroupRouts.DELETE("/:student_group_id/", s.studentGroupController.Delete)

	lessonRouts := server.Group("/lesson")
	lessonRouts.GET("/", s.lessonController.GetAll)
	lessonRouts.POST("/", s.lessonController.Create)
	lessonRouts.PUT("/:lesson_id/", s.lessonController.Update)
	lessonRouts.DELETE("/:lesson_id/", s.lessonController.Delete)
	lessonRouts.POST("/swap/", s.lessonController.SwapSlots)

	err := server.Run(s.listenAddr)
	return err
}
