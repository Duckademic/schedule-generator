package services

import (
	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

// LessonTypeService aggregates and manages lesson types that the generator works with.
type LessonTypeService interface {
	Find(uuid.UUID) *entities.LessonType // Returns a pointer to the lesson type with the given ID
	GetAll() []*entities.LessonType      // Returns an array with all lesson types as pointers
}

// NewLessonTypeService creates new LessonTypeService instance.
//
// It requires an array of database lesson types (lTypes).
func NewLessonTypeService(lTypes []types.LessonType) (LessonTypeService, error) {
	lts := lessonTypeService{
		lessonTypes: make([]*entities.LessonType, len(lTypes)),
	}

	for i, lt := range lTypes {
		lts.lessonTypes[i] = &entities.LessonType{
			ID:    lt.ID,
			Name:  lt.Name,
			Weeks: lt.Weeks,
			Value: lt.Value,
		}
	}

	return &lts, nil
}

type lessonTypeService struct {
	lessonTypes []*entities.LessonType
}

func (lts *lessonTypeService) Find(id uuid.UUID) *entities.LessonType {
	for i := range lts.lessonTypes {
		if lts.lessonTypes[i].ID == id {
			return lts.lessonTypes[i]
		}
	}
	return nil
}

func (lts *lessonTypeService) GetAll() []*entities.LessonType {
	return lts.lessonTypes
}
