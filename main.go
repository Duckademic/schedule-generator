package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Duckademic/schedule-generator/generator"
	"github.com/Duckademic/schedule-generator/repositories"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func init() {
	if err := ENVLoad(); err != nil {
		log.Fatal("Init error: " + err.Error())
	}

	db, err := repositories.InitDB()
	if err != nil {
		log.Fatal("Init error: " + err.Error())
	}

	s, err := ServerInit(db)
	if err != nil {
		log.Fatal("Init error: " + err.Error())
	}
	server = *s
}

var server JSONAPIServer

func main() {
	wl := [][]float32{
		{},
		{0.6, 2, 1.8, 1.6, 1.4, 1.2, 1.0},
		{0.6, 2, 1.8, 1.6, 1.4, 1.2, 1.0},
		{0.6, 2, 1.8, 1.6, 1.4, 1.2, 1.0},
		{0.6, 2, 1.8, 1.6, 1.4, 1.2, 1.0},
		{0.6, 2, 1.8, 1.6, 1.4, 1.2, 1.0},
		{},
	}

	var gen *generator.ScheduleGenerator
	func() {
		log.Println("GENERATING start")
		defer func(start time.Time) {
			log.Println("GENERATING finished " + time.Since(start).String())
		}(time.Now())

		var err error
		gen, err = generator.NewScheduleGenerator(generator.ScheduleGeneratorConfig{
			LessonsValue:       2,
			Start:              time.Date(2025, time.January, 19, 0, 0, 0, 0, time.UTC),
			End:                time.Date(2025, time.April, 30, 0, 0, 0, 0, time.UTC),
			WorkLessons:        wl,
			MaxStudentWorkload: 4,
			FillPercentage:     1,
		})
		if err != nil {
			panic(err)
		}

		teachers, sGroups, disciplines, sLoads, lTypes := INIT()
		err = gen.SetStudentGroups(sGroups)
		if err != nil {
			panic(err)
		}
		err = gen.SetTeachers(teachers)
		if err != nil {
			panic(err)
		}
		err = gen.SetDisciplines(disciplines)
		if err != nil {
			panic(err)
		}
		err = gen.SetLessonTypes(lTypes)
		if err != nil {
			panic(err)
		}
		err = gen.SetStudyLoads(sLoads)
		if err != nil {
			panic(err)
		}

		err = gen.GenerateSchedule()
		if err != nil {
			panic(err)
		}
	}()
	gen.CheckSchedule()
	gen.WriteSchedule()

	err := server.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ENVLoad() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %s", err.Error())
	}

	return nil
}

func ServerInit(db *gorm.DB) (*JSONAPIServer, error) {
	release := os.Getenv("RELEASE")
	if release == "1" {
		gin.SetMode(gin.ReleaseMode)
	}

	wl := [][]float32{
		{},
		{0.6, 2, 1.8, 1.6, 1.4, 1.2, 1.0},
		{0.6, 2, 1.8, 1.6, 1.4, 1.2, 1.0},
		{0.6, 2, 1.8, 1.6, 1.4, 1.2, 1.0},
		{0.6, 2, 1.8, 1.6, 1.4, 1.2, 1.0},
		{0.6, 2, 1.8, 1.6, 1.4, 1.2, 1.0},
		{},
	}
	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("port not specified at .env file")
	}

	server, err := NewJSONAPIServer(fmt.Sprintf("localhost:%s", port), generator.ScheduleGeneratorConfig{
		LessonsValue:       2,
		Start:              time.Date(2025, time.January, 19, 0, 0, 0, 0, time.UTC),
		End:                time.Date(2025, time.May, 31, 0, 0, 0, 0, time.UTC),
		WorkLessons:        wl,
		MaxStudentWorkload: 4,
	}, db)

	if err != nil {
		return nil, fmt.Errorf("server creation error: %s", err.Error())
	}

	return server, nil
}
