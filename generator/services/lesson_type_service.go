package services

import (
	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type LessonTypeService interface {
	Find(uuid.UUID) *entities.LessonType
	GetAll() []entities.LessonType
	GetWeekOffset() int
}

func NewLessonTypeService(lTypes []types.LessonType) (LessonTypeService, error) {
	lts := lessonTypeService{
		lessonTypes: make([]entities.LessonType, len(lTypes)),
	}

	for i, lt := range lTypes {
		lts.lessonTypes[i] = entities.LessonType{
			ID:    lt.ID,
			Name:  lt.Name,
			Weeks: lt.Weeks,
			Value: lt.Value,
		}
	}

	return &lts, nil
}

type lessonTypeService struct {
	lessonTypes []entities.LessonType
}

func (lts *lessonTypeService) Find(id uuid.UUID) *entities.LessonType {
	for i := range lts.lessonTypes {
		if lts.lessonTypes[i].ID == id {
			return &lts.lessonTypes[i]
		}
	}
	return nil
}

func (lts *lessonTypeService) GetAll() []entities.LessonType {
	return lts.lessonTypes
}

func (lts *lessonTypeService) GetWeekOffset() (maxW int) {
	for _, lType := range lts.lessonTypes {
		maxW = max(maxW, lType.Weeks)
	}

	return
}
