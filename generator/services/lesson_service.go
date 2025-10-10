package services

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/generator/entities"
)

type LessonService interface {
	GetAll() []entities.Lesson
	AddLesson(*entities.Teacher, *entities.StudentGroup, *entities.Discipline, entities.LessonSlot, *entities.LessonType) error
	GetWeekLessons(int) []entities.Lesson
}

func NewLessonService(lessonValue int) (LessonService, error) {
	if lessonValue <= 0 {
		return nil, fmt.Errorf("lessonValue under/equal 0 (%d)", lessonValue)
	}

	ls := lessonService{lessonValue: lessonValue}

	return &ls, nil
}

type lessonService struct {
	lessons     []entities.Lesson
	lessonValue int
}

func (ls *lessonService) GetAll() []entities.Lesson {
	return ls.lessons
}

func (ls *lessonService) AddLesson(
	teacher *entities.Teacher,
	studentGroup *entities.StudentGroup,
	discipline *entities.Discipline,
	slot entities.LessonSlot,
	lType *entities.LessonType,
) error {
	// загальні перевірки
	if teacher == nil {
		return fmt.Errorf("teacher can't be nil")
	}
	if studentGroup == nil {
		return fmt.Errorf("student group can't be nil")
	}
	if discipline == nil {
		return fmt.Errorf("discipline can't be nil")
	}

	lesson := &entities.Lesson{
		Teacher:      teacher,
		StudentGroup: studentGroup,
		Discipline:   discipline,
		Slot:         slot,
		Type:         lType,
		Value:        ls.lessonValue,
	}

	if err := teacher.CheckLesson(lesson); err != nil {
		return err
	}
	if err := studentGroup.CheckLesson(lesson); err != nil {
		return err
	}

	// перевірки дисципліни
	if discipline.EnoughHours() {
		return fmt.Errorf("discipline have enough hours")
	}

	// додавання пари
	ls.lessons = append(ls.lessons, *lesson)

	lesson.StudentGroup.AddLesson(lesson, true)
	lesson.Teacher.AddLesson(lesson, true)

	lesson.Discipline.Load[0].CurrentHours += ls.lessonValue
	return nil
}

func (ls *lessonService) GetWeekLessons(week int) (res []entities.Lesson) {
	for _, l := range ls.lessons {
		if l.Slot.Day/7 == week {
			res = append(res, l)
		}
	}
	return
}
