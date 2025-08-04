package generator

import (
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

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

func (d *Discipline) CountHourDeficit() (count int) {
	for _, load := range d.Load {
		count += load.CountHourDeficit()
	}

	return
}

type DisciplineService interface {
	GetAll() []Discipline
	Find(uuid.UUID) *Discipline
	CountHourDeficit() int
}

func NewDisciplineService(disciplines []types.Discipline) (DisciplineService, error) {
	ds := disciplineService{disciplines: make([]Discipline, len(disciplines))}

	for i := range disciplines {
		ds.disciplines[i] = Discipline{
			ID:   disciplines[i].ID,
			Name: disciplines[i].Name,
		}
	}

	return &ds, nil
}

type disciplineService struct {
	disciplines []Discipline
}

func (ds *disciplineService) GetAll() []Discipline {
	return ds.disciplines
}

func (ds *disciplineService) Find(disciplineID uuid.UUID) *Discipline {
	for i := range ds.disciplines {
		if ds.disciplines[i].ID == disciplineID {
			return &ds.disciplines[i]
		}
	}

	return nil
}

func (ds *disciplineService) CountHourDeficit() (count int) {
	for _, d := range ds.disciplines {
		count += d.CountHourDeficit()
	}
	return
}
