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
	Type         *LessonType
	Teacher      *Teacher
	StudentGroup *StudentGroup
	Discipline   *Discipline
}

type LessonService interface {
	GetAll() []Lesson
	CreateWithoutChecks(*Teacher, *StudentGroup, *Discipline, LessonSlot, *LessonType)
	CreateWithChecks(*Teacher, *StudentGroup, *Discipline, LessonSlot, *LessonType) error
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

func (ls *lessonService) CreateWithoutChecks(
	teacher *Teacher,
	studentGroup *StudentGroup,
	discipline *Discipline,
	slot LessonSlot,
	lType *LessonType,
) {
	ls.lessons = append(ls.lessons, Lesson{
		Teacher:      teacher,
		StudentGroup: studentGroup,
		Discipline:   discipline,
		Slot:         slot,
		Type:         lType,
	})

	teacher.SetOneSlotBusyness(slot, true)
	studentGroup.SetOneSlotBusyness(slot, true)
	discipline.CurrentHours += ls.lessonValue
}

func (ls *lessonService) CreateWithChecks(
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

	// перевірки викладача
	if err := teacher.CheckSlot(slot); err != nil {
		return err
	}
	if teacher.IsBusy(slot) {
		return fmt.Errorf("teacher is busy")
	}

	// перевірки групи студентів
	if err := teacher.CheckSlot(slot); err != nil {
		return err
	}
	if studentGroup.IsBusy(slot) {
		return fmt.Errorf("student group is busy")
	}

	// перевірки дисципліни
	if discipline.EnoughHours() {
		return fmt.Errorf("discipline have enough hours")
	}

	ls.CreateWithoutChecks(teacher, studentGroup, discipline, slot, lType)
	return nil
}

func (ls *lessonService) GetWeekLessons(week int) (res []Lesson) {
	for _, l := range ls.lessons {
		if l.Slot.Day/7 == week {
			res = append(res, l)
		}
	}
	return
}
