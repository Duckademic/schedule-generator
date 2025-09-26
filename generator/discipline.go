package generator

import "github.com/google/uuid"

type DisciplineLoad struct {
	Teacher    *Teacher
	Groups     []*StudentGroup
	LessonType *LessonType
	LessonChecker
}

type Discipline struct {
	ID   uuid.UUID
	Name string
	Load []DisciplineLoad
}

func (d *Discipline) AddLoad(teacher *Teacher, hours int, groups []*StudentGroup, lType *LessonType) error {
	dl := DisciplineLoad{
		LessonChecker: LessonChecker{
			RequiredHours: hours * len(groups),
		},
		Teacher:    teacher,
		Groups:     groups,
		LessonType: lType,
	}

	d.Load = append(d.Load, dl)
	return nil
}

// ПЕРЕПИСАТИ
func (d *Discipline) EnoughHours() bool {
	return d.Load[0].RequiredHours <= d.Load[0].CurrentHours
}

// Returns sum of pending hours of all disciplines.
// Time complexity O(n)
func (d *Discipline) CountHourDeficit() (count int) {
	for _, load := range d.Load {
		count += load.CountHourDeficit()
	}

	return
}
