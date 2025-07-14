package services

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type StudentGroupServise interface {
	Service[types.StudentGroup]
}

type studentGroupServise struct {
	studentGroups []types.StudentGroup
}

func NewStudentGroupService(studentGroups []types.StudentGroup) StudentGroupServise {
	sgs := studentGroupServise{studentGroups: studentGroups}

	return &sgs
}

func (sgs *studentGroupServise) Create(group types.StudentGroup) (*types.StudentGroup, error) {
	if sgs.Find(group.ID) != nil {
		return nil, fmt.Errorf("student group %s already exists", group.ID.String())
	}

	sgs.studentGroups = append(sgs.studentGroups, group)
	return &group, nil
}

func (sgs *studentGroupServise) Update(group types.StudentGroup) error {
	g := sgs.Find(group.ID)
	if g == nil {
		return fmt.Errorf("student group %s not found", group.ID.String())
	}

	g.Name = group.Name
	return nil
}

func (sgs *studentGroupServise) Delete(groupId uuid.UUID) error {
	for i, group := range sgs.studentGroups {
		if group.ID == groupId {
			sgs.studentGroups = append(sgs.studentGroups[:i], sgs.studentGroups[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("student group %s not found", groupId.String())
}

func (sgs *studentGroupServise) GetAll() []types.StudentGroup {
	return sgs.studentGroups
}

// return will be nil if not found
func (sgs *studentGroupServise) Find(id uuid.UUID) *types.StudentGroup {
	var group *types.StudentGroup
	for i := range sgs.studentGroups {
		if sgs.studentGroups[i].ID == id {
			group = &sgs.studentGroups[i]
			break
		}
	}

	return group
}
