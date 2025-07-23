package services

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type StudentGroupService interface {
	Service[types.StudentGroup]
}

type studentGroupService struct {
	studentGroups []types.StudentGroup
}

func NewStudentGroupService(studentGroups []types.StudentGroup) StudentGroupService {
	sgs := studentGroupService{studentGroups: studentGroups}

	return &sgs
}

func (sgs *studentGroupService) Create(group types.StudentGroup) (*types.StudentGroup, error) {
	if sgs.Find(group.ID) != nil {
		return nil, fmt.Errorf("student group %s already exists", group.ID.String())
	}

	sgs.studentGroups = append(sgs.studentGroups, group)
	return &group, nil
}

func (sgs *studentGroupService) Update(group types.StudentGroup) error {
	g := sgs.Find(group.ID)
	if g == nil {
		return fmt.Errorf("student group %s not found", group.ID.String())
	}

	g.Name = group.Name
	return nil
}

func (sgs *studentGroupService) Delete(groupId uuid.UUID) error {
	for i, group := range sgs.studentGroups {
		if group.ID == groupId {
			sgs.studentGroups = append(sgs.studentGroups[:i], sgs.studentGroups[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("student group %s not found", groupId.String())
}

func (sgs *studentGroupService) GetAll() []types.StudentGroup {
	return sgs.studentGroups
}

// return will be nil if not found
func (sgs *studentGroupService) Find(id uuid.UUID) *types.StudentGroup {
	var group *types.StudentGroup
	for i := range sgs.studentGroups {
		if sgs.studentGroups[i].ID == id {
			group = &sgs.studentGroups[i]
			break
		}
	}

	return group
}
