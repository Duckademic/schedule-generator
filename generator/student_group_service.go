package generator

import (
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type StudentGroupService interface {
	GetAll() []StudentGroup
	Find(uuid.UUID) *StudentGroup
	CountWindows() int
	CountHourDeficit() int
	CountLessonOverlapping() int
}

type studentGroupService struct {
	studentGroups []StudentGroup
}

func NewStudentGroupService(studentGroups []types.StudentGroup, maxLessonsPerDay int, busyGrid [][]float32) (StudentGroupService, error) {
	sgs := studentGroupService{
		studentGroups: make([]StudentGroup, len(studentGroups)),
	}

	for i := range studentGroups {
		sgs.studentGroups[i] = StudentGroup{
			ID:                studentGroups[i].ID,
			Name:              studentGroups[i].Name,
			MaxLessonsPerDay:  maxLessonsPerDay,
			LessonTypeBinding: map[*LessonType]*StudentGroupLoad{},
		}
		studentGroup := &sgs.studentGroups[i]
		studentGroup.BusyGrid = *NewBusyGrid(busyGrid)

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

func (sgs *studentGroupService) GetAll() []StudentGroup {
	return sgs.studentGroups
}

// return nil if not found
func (sgs *studentGroupService) Find(id uuid.UUID) *StudentGroup {
	for i := range sgs.studentGroups {
		if sgs.studentGroups[i].ID == id {
			return &sgs.studentGroups[i]
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
