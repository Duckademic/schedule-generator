package generator

import (
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type LessonType struct {
	ID    uuid.UUID
	Name  string
	Weeks int
}

type LessonTypeService interface {
	Find(uuid.UUID) *LessonType
	GetAll() []LessonType
	GetWeekOffset() int
}

func NewLessonTypeService(lTypes []types.LessonType) (LessonTypeService, error) {
	lts := lessonTypeService{
		lessonTypes: make([]LessonType, len(lTypes)),
	}

	for i, lt := range lTypes {
		lts.lessonTypes[i] = LessonType{
			ID:    lt.ID,
			Name:  lt.Name,
			Weeks: lt.Weeks,
		}
	}

	return &lts, nil
}

type lessonTypeService struct {
	lessonTypes []LessonType
}

func (lts *lessonTypeService) Find(id uuid.UUID) *LessonType {
	for i := range lts.lessonTypes {
		if lts.lessonTypes[i].ID == id {
			return &lts.lessonTypes[i]
		}
	}
	return nil
}

func (lts *lessonTypeService) GetAll() []LessonType {
	return lts.lessonTypes
}

func (lts *lessonTypeService) GetWeekOffset() (maxW int) {
	for _, lType := range lts.lessonTypes {
		maxW = max(maxW, lType.Weeks)
	}

	return
}
