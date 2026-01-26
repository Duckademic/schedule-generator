package services

import (
	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

// DisciplineService aggregates and manages disciplines that the generator works with.
type DisciplineService interface {
	Find(uuid.UUID) *entities.Discipline // Returns a pointer to the discipline with the given ID.
	GetAll() []*entities.Discipline      // Returns an array with all disciplines as pointers.
}

// NewDisciplineService creates a new DisciplineService basic instance.
//
// It requires an array of database disciplines (d).
func NewDisciplineService(d []types.Discipline) (DisciplineService, error) {
	ds := disciplineService{disciplines: make([]*entities.Discipline, len(d))}

	for i := range d {
		ds.disciplines[i] = &entities.Discipline{
			ID:   d[i].ID,
			Name: d[i].Name,
		}
	}

	return &ds, nil
}

// disciplineService is the basic implementation of the DisciplineService interface.
type disciplineService struct {
	disciplines []*entities.Discipline
}

func (ds *disciplineService) GetAll() []*entities.Discipline {
	return ds.disciplines
}
func (ds *disciplineService) Find(disciplineID uuid.UUID) *entities.Discipline {
	for i := range ds.disciplines {
		if ds.disciplines[i].ID == disciplineID {
			return ds.disciplines[i]
		}
	}

	return nil
}
