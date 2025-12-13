package services

import (
	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type StudentGroupService interface {
	GetAll() []*entities.StudentGroup
	Find(uuid.UUID) *entities.StudentGroup
	CountWindows() int // Returns sum of all student groups windows
	CountHourDeficit() int
	CountLessonOverlapping() int   // Returns sum of all lesson overlap.
	CountOvertimeLessons() int     // Return sum of all overtime lesson (above the daily limit).
	CountInvalidLessonsTypes() int // Returns sum of all lesson scheduled on days that are not allowed for their types.
}

type studentGroupService struct {
	studentGroups []*entities.StudentGroup
}

func NewStudentGroupService(studentGroups []types.StudentGroup, maxLessonsPerDay int, busyGrid [][]float32) (StudentGroupService, error) {
	sgs := studentGroupService{
		studentGroups: make([]*entities.StudentGroup, len(studentGroups)),
	}

	for i := range studentGroups {
		sgs.studentGroups[i] = &entities.StudentGroup{
			ID:                studentGroups[i].ID,
			Name:              studentGroups[i].Name,
			MaxLessonsPerDay:  maxLessonsPerDay,
			LessonTypeBinding: map[*entities.LessonType]*entities.StudentGroupLoad{},
		}
		studentGroup := sgs.studentGroups[i]
		studentGroup.BusyGrid = *entities.NewBusyGrid(busyGrid)

		md := studentGroups[i].MilitaryDay - 1
		if md != -1 {
			if err := studentGroup.CheckWeekDay(md); err != nil {
				return nil, err
			}
			studentGroup.SetDayBusyness(make([]float32, len(studentGroup.Grid[md])), md)
		}
	}

	return &sgs, nil
}

func (sgs *studentGroupService) GetAll() []*entities.StudentGroup {
	return sgs.studentGroups
}

// return nil if not found
func (sgs *studentGroupService) Find(id uuid.UUID) *entities.StudentGroup {
	for i := range sgs.studentGroups {
		if sgs.studentGroups[i].ID == id {
			return sgs.studentGroups[i]
		}
	}

	return nil
}

func (sgs *studentGroupService) CountWindows() (count int) {
	for _, g := range sgs.studentGroups {
		count += g.CountWindows()
	}
	return
}

func (sgs *studentGroupService) CountHourDeficit() (count int) {
	for _, studentGroup := range sgs.studentGroups {
		count += studentGroup.CountHourDeficit()
	}

	return count
}

func (sgs *studentGroupService) CountLessonOverlapping() (count int) {
	for _, studentGroup := range sgs.studentGroups {
		count += studentGroup.CountLessonOverlapping()
	}

	return
}

func (sgs *studentGroupService) CountOvertimeLessons() (count int) {
	for _, sg := range sgs.studentGroups {
		count += sg.GetOvertimeLessons()
	}
	return
}

func (sgs *studentGroupService) CountInvalidLessonsTypes() (count int) {
	for _, sg := range sgs.studentGroups {
		count += sg.GetInvalidLessonsType()
	}
	return
}
