package generator

import "fmt"

type LessonService interface {
	GetAll() []Lesson
	AddWithoutChecks(*Teacher, *StudentGroup, *Discipline, LessonSlot, *LessonType)
	AddWithChecks(*Teacher, *StudentGroup, *Discipline, LessonSlot, *LessonType) error
	GetWeekLessons(int) []Lesson
}

func NewLessonService(lessonValue int) (LessonService, error) {
	if lessonValue <= 0 {
		return nil, fmt.Errorf("lessonValue under/equal 0 (%d)", lessonValue)
	}

	ls := lessonService{lessonValue: lessonValue}

	return &ls, nil
}

type lessonService struct {
	lessons     []Lesson
	lessonValue int
}

func (ls *lessonService) GetAll() []Lesson {
	return ls.lessons
}

func (ls *lessonService) AddWithoutChecks(
	teacher *Teacher,
	studentGroup *StudentGroup,
	discipline *Discipline,
	slot LessonSlot,
	lType *LessonType,
) {
	ls.AddLesson(ls.CreateLesson(teacher, studentGroup, discipline, slot, lType))
}

func (ls *lessonService) AddWithChecks(
	teacher *Teacher,
	studentGroup *StudentGroup,
	discipline *Discipline,
	slot LessonSlot,
	lType *LessonType,
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

	lesson := ls.CreateLesson(teacher, studentGroup, discipline, slot, lType)

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

	ls.AddLesson(lesson)
	return nil
}

func (ls *lessonService) AddLesson(l *Lesson) {
	ls.lessons = append(ls.lessons, *l)

	l.StudentGroup.AddLesson(l, true)
	l.Teacher.AddLesson(l, true)

	l.Discipline.Load[0].CurrentHours += ls.lessonValue
}

func (ls *lessonService) CreateLesson(
	teacher *Teacher,
	studentGroup *StudentGroup,
	discipline *Discipline,
	slot LessonSlot,
	lType *LessonType,
) *Lesson {
	return &Lesson{
		Teacher:      teacher,
		StudentGroup: studentGroup,
		Discipline:   discipline,
		Slot:         slot,
		Type:         lType,
		Value:        ls.lessonValue,
	}
}

func (ls *lessonService) GetWeekLessons(week int) (res []Lesson) {
	for _, l := range ls.lessons {
		if l.Slot.Day/7 == week {
			res = append(res, l)
		}
	}
	return
}
