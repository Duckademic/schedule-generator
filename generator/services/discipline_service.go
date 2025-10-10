package services

import (
	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type DisciplineService interface {
	GetAll() []entities.Discipline
	Find(uuid.UUID) *entities.Discipline
	CountHourDeficit() int
}

func NewDisciplineService(disciplines []types.Discipline) (DisciplineService, error) {
	ds := disciplineService{disciplines: make([]entities.Discipline, len(disciplines))}

	for i := range disciplines {
		ds.disciplines[i] = entities.Discipline{
			ID:   disciplines[i].ID,
			Name: disciplines[i].Name,
		}
	}

	return &ds, nil
}

type disciplineService struct {
	disciplines []entities.Discipline
}

func (ds *disciplineService) GetAll() []entities.Discipline {
	return ds.disciplines
}

func (ds *disciplineService) Find(disciplineID uuid.UUID) *entities.Discipline {
	for i := range ds.disciplines {
		if ds.disciplines[i].ID == disciplineID {
			return &ds.disciplines[i]
		}
	}

	return nil
}

// Returns hour deficit of all disciplines.
// Time complexity O(n)
func (ds *disciplineService) CountHourDeficit() (count int) {
	for _, d := range ds.disciplines {
		count += d.CountHourDeficit()
	}
	return
}
