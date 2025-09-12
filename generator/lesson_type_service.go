package generator

import (
	"github.com/google/uuid"
)

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
