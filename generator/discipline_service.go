package generator

import (
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type DisciplineService interface {
	GetAll() []Discipline
	Find(uuid.UUID) *Discipline
}

type Discipline struct {
	ID           uuid.UUID
	Name         string
	LoadHours    int
	CurrentHours int
	// Lessons map[string]int // тип - кількість годин
}

func (d *Discipline) EnoughHours() bool {
	return d.LoadHours <= d.CurrentHours
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
