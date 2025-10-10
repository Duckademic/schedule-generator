package services

import (
	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type StudentGroupService interface {
	GetAll() []entities.StudentGroup
	Find(uuid.UUID) *entities.StudentGroup
	CountWindows() int
	CountHourDeficit() int
	CountLessonOverlapping() int
}

type studentGroupService struct {
	studentGroups []entities.StudentGroup
}

func NewStudentGroupService(studentGroups []types.StudentGroup, maxLessonsPerDay int, busyGrid [][]float32) (StudentGroupService, error) {
	sgs := studentGroupService{
		studentGroups: make([]entities.StudentGroup, len(studentGroups)),
	}

	for i := range studentGroups {
		sgs.studentGroups[i] = entities.StudentGroup{
			ID:                studentGroups[i].ID,
			Name:              studentGroups[i].Name,
			MaxLessonsPerDay:  maxLessonsPerDay,
			LessonTypeBinding: map[*entities.LessonType]*entities.StudentGroupLoad{},
		}
		studentGroup := &sgs.studentGroups[i]
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

func (sgs *studentGroupService) GetAll() []entities.StudentGroup {
	return sgs.studentGroups
}

// return nil if not found
func (sgs *studentGroupService) Find(id uuid.UUID) *entities.StudentGroup {
	for i := range sgs.studentGroups {
		if sgs.studentGroups[i].ID == id {
			return &sgs.studentGroups[i]
		}
	}

	return nil
}

// Returns sum of all student groups windows
// Time complexity O(n)
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

// Returns sum of all lesson overlap.
// Time complexity O(n^2)
func (sgs *studentGroupService) CountLessonOverlapping() (count int) {
	for _, studentGroup := range sgs.studentGroups {
		count += studentGroup.CountLessonOverlapping()
	}

	return
}
