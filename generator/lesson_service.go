package generator

import "fmt"

type LessonType struct {
	Name string
}

type LessonSlot struct {
	Day  int
	Slot int
}

type Lesson struct {
	// ID           uuid.UUID
	Slot         LessonSlot
	Value        int // кількість академічних годин
	Type         LessonType
	Teacher      *Teacher
	StudentGroup *StudentGroup
	Discipline   *Discipline
}

type LessonService interface {
	GetAll() []Lesson
	CreateWithoutChecks(*Teacher, *StudentGroup, *Discipline, LessonSlot, LessonType)
	CreateWithChecks(*Teacher, *StudentGroup, *Discipline, LessonSlot, LessonType) error
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

func (ls *lessonService) CreateWithoutChecks(
	teacher *Teacher,
	studentGroup *StudentGroup,
	discipline *Discipline,
	slot LessonSlot,
	lType LessonType,
) {
	ls.lessons = append(ls.lessons, Lesson{
		Teacher:      teacher,
		StudentGroup: studentGroup,
		Discipline:   discipline,
		Slot:         slot,
		Type:         lType,
	})
}

func (ls *lessonService) CreateWithChecks(
	teacher *Teacher,
	studentGroup *StudentGroup,
	discipline *Discipline,
	slot LessonSlot,
	lType LessonType,
) error {
	if teacher == nil {
		return fmt.Errorf("teacher can't be nil")
	}
	if studentGroup == nil {
		return fmt.Errorf("student group can't be nil")
	}
	if discipline == nil {
		return fmt.Errorf("discipline can't be nil")
	}
	// дописати перевірки

	ls.CreateWithoutChecks(teacher, studentGroup, discipline, slot, lType)
	return nil
}
