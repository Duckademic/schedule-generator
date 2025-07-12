package main

import (
	"flag"
	"log"
	"time"

	"github.com/Duckademic/schedule-generator/generator"
)

func main() {
	gen, err := generator.NewScheduleGenerator(generator.ScheduleGeneratorConfig{
		LessonsValue:       2,
		Start:              time.Date(2025, time.January, 19, 0, 0, 0, 0, time.UTC),
		End:                time.Date(2025, time.May, 30, 0, 0, 0, 0, time.UTC),
		WorkLessons:        []int{0, 7, 7, 7, 7, 7, 0},
		MaxStudentWorkload: 3,
	})
	if err != nil {
		log.Fatal("invalid generator config: " + err.Error())
	}

	gen.GenerateShedule(INIT())

	listenAddr := flag.String("listenaddr", ":8080", "listen address the service is running")
	flag.Parse()

	server, err := NewJSONAPIServer(*listenAddr, generator.ScheduleGeneratorConfig{
		LessonsValue:       2,
		Start:              time.Date(2025, time.January, 19, 0, 0, 0, 0, time.UTC),
		End:                time.Date(2025, time.May, 30, 0, 0, 0, 0, time.UTC),
		WorkLessons:        []int{0, 7, 7, 7, 7, 7, 0},
		MaxStudentWorkload: 3,
	})

	if err != nil {
		log.Fatal("Server creation error: " + err.Error())
	}

	server.Run()
}
