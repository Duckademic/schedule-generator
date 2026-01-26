package services

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

// StudentGroupService aggregates and manages student groups that the generator works with.
type StudentGroupService interface {
	Find(uuid.UUID) *entities.StudentGroup // Returns a pointer to the student group with the given ID.
	GetAll() []*entities.StudentGroup      // Returns a slice with all student groups as pointers.
	CountWindows() int                     // Returns the sum of windows (gaps between busy slots).
	CountHourDeficit() int                 // Returns the number of missing study hours.
	CountLessonOverlapping() int           // Returns the count of overlapping lessons.
	CountOvertimeLessons() int             // Returns the total number of overtime lessons (above the daily limit).
	// Returns the total number of lesson scheduled on days that are not allowed for their type.
	CountInvalidLessonsByType() int
}

// NewStudentGroupService creates a new StudentGroupService basic instance.
//
// It requires an array of database student groups (sg), day load limit (dll), and a busy grid for them (bg).
//
// Returns an error if any student group is an invalid model.
func NewStudentGroupService(sg []types.StudentGroup, dl int, bg [][]float32) (StudentGroupService, error) {
	sgs := studentGroupService{
		studentGroups: make([]*entities.StudentGroup, len(sg)),
	}

	for i := range sg {
		sgs.studentGroups[i] = entities.NewDefaultStudentGroup(sg[i].ID, sg[i].Name, dl, entities.NewBusyGrid(bg))
		studentGroup := sgs.studentGroups[i]

		// set military day by marks slots on this day as blocked
		md := sg[i].MilitaryDay
		if md != -1 {
			if err := studentGroup.CheckWeekDay(md); err != nil {
				return nil, err
			}
			studentGroup.BlockWeekDay(md)
		}
	}

	// create connection for student groups
	for i := range sg {
		mainGroup := sgs.Find(sg[i].ID)
		for _, gID := range sg[i].ConnectedGroups {
			otherG := sgs.Find(gID)
			if otherG == nil {
				return nil, fmt.Errorf("Can't find connected group %s for group %s (%s)", gID, sg[i].Name, sg[i].ID)
			}
			otherG.AddConnectedGroup(mainGroup)
		}
	}

	return &sgs, nil
}

// studentGroupService is the basic implementation of the StudentGroupService interface.
type studentGroupService struct {
	studentGroups []*entities.StudentGroup
}

func (sgs *studentGroupService) GetAll() []*entities.StudentGroup {
	return sgs.studentGroups
}
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
		count += studentGroup.CountLessonOverlapping(studentGroup.GetAssignedLessons())
	}

	return
}
func (sgs *studentGroupService) CountOvertimeLessons() (count int) {
	for _, sg := range sgs.studentGroups {
		count += sg.CountOvertimeLessons()
	}
	return
}
func (sgs *studentGroupService) CountInvalidLessonsByType() (count int) {
	for _, sg := range sgs.studentGroups {
		count += sg.CountInvalidLessonsByType()
	}
	return
}
