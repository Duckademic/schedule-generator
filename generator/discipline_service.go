package generator

import (
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

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
