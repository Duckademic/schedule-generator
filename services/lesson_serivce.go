package services

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type LessonService interface {
	Service[types.Lesson]
	SwapSlots(firstId, secondId uuid.UUID) error
}

type lessonService struct {
	lessons []types.Lesson
}

func NewLessonService(lessons []types.Lesson) LessonService {
	ls := lessonService{
		lessons: lessons,
	}

	return &ls
}

func (ls *lessonService) Create(lesson types.Lesson) (*types.Lesson, error) {
	lesson.ID = uuid.New()
	ls.lessons = append(ls.lessons, lesson)

	return &lesson, nil
}

func (ls *lessonService) Update(lesson types.Lesson) error {
	l := ls.Find(lesson.ID)
	if l == nil {
		return fmt.Errorf("lesson %s not found", lesson.ID)
	}

	l.StartTime = lesson.StartTime
	l.EndTime = lesson.EndTime
	l.Value = lesson.Value
	l.Type = lesson.Type
	return nil
}

func (ls *lessonService) Find(lessonId uuid.UUID) *types.Lesson {
	for i := range ls.lessons {
		if ls.lessons[i].ID == lessonId {
			return &ls.lessons[i]
		}
	}
	return nil
}

func (ls *lessonService) Delete(lessonId uuid.UUID) error {
	for i := range ls.lessons {
		if ls.lessons[i].ID == lessonId {
			ls.lessons = append(ls.lessons[:i], ls.lessons[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("lesson %s not found", lessonId)
}

func (ls *lessonService) GetAll() []types.Lesson {
	return ls.lessons
}

func (ls *lessonService) SwapSlots(firstId, secondId uuid.UUID) error {
	first := ls.Find(firstId)
	if first == nil {
		return fmt.Errorf("lesson %s not found (first)", firstId)
	}

	second := ls.Find(secondId)
	if second == nil {
		return fmt.Errorf("lesson %s not found (second)", firstId)
	}

	tmpTime := first.StartTime
	first.StartTime = second.StartTime
	second.StartTime = tmpTime

	tmpTime = first.EndTime
	first.EndTime = second.EndTime
	second.EndTime = tmpTime
	return nil
}
