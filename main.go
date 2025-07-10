package main

import (
	"log"
	"time"

	"github.com/Duckademic/schedule-generator/generator"
	"github.com/Duckademic/schedule-generator/services"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/gin-gonic/gin"
)

var ts = services.SimpleService[*types.Teacher]{}
var sgs = services.SimpleService[*types.StudentGroup]{}
var ds = services.SimpleService[*types.Discipline]{}
var sls = services.SimpleService[*types.StudyLoad]{}

func main() {
	// types.TestCheckWindows()

	INIT()

	log.Printf("%d %d %d %d", len(ts.GetAll()), len(sgs.GetAll()), len(ds.GetAll()), len(sls.GetAll()))

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

	// Convert []*types.Teacher to []types.Teacher
	teachersPtr := ts.GetAll()
	teachers := make([]types.Teacher, len(teachersPtr))
	for i, t := range teachersPtr {
		if t != nil {
			teachers[i] = *t
		}
	}

	teacherService, err := services.NewTeacherService(teachers)
	if err != nil {
		log.Fatal("failed to create TeacherService: " + err.Error())
	}

	// Convert []*types.StudentGroup to []types.StudentGroup
	sgsPtr := sgs.GetAll()
	studentGroups := make([]types.StudentGroup, len(sgsPtr))
	for i, sg := range sgsPtr {
		if sg != nil {
			studentGroups[i] = *sg
		}
	}

	studentGroupServise, err := services.NewStudentGroupService(studentGroups, 4)
	if err != nil {
		log.Fatal("failed to create StudentGroupService: " + err.Error())
	}

	// Convert []*types.StudyLoad to []types.StudyLoad
	slsPtr := sls.GetAll()
	studyLoads := make([]types.StudyLoad, len(slsPtr))
	for i, sl := range slsPtr {
		if sl != nil {
			studyLoads[i] = *sl
		}
	}

	defer func(start time.Time) {
		log.Printf("Розклад згенеровано за %v", time.Since(start))
	}(time.Now())
	gen.GenerateShadule(studentGroupServise, teacherService, studyLoads)

	return

	server := gin.Default()
	server.Run(":8080")
}
