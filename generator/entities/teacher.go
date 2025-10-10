package entities

import (
	"fmt"

	"github.com/google/uuid"
)

type TeacherLoad struct {
	Discipline *Discipline
	Groups     []*StudentGroup
	LessonType *LessonType
}

type Teacher struct {
	BusyGrid
	LessonChecker
	ID       uuid.UUID
	UserName string
	Priority int
	Load     []TeacherLoad
}

func (t *Teacher) AddLoad(discipline *Discipline, lessonType *LessonType, groups []*StudentGroup, hours int) error {
	tl := TeacherLoad{
		Discipline: discipline,
		LessonType: lessonType,
		Groups:     groups,
	}

	t.RequiredHours += hours

	t.Load = append(t.Load, tl)
	return nil
}

func (t *Teacher) AddLesson(lesson *Lesson, ignoreCheck bool) error {
	err := t.CheckLesson(lesson)
	if err != nil && !ignoreCheck {
		return err
	}

	t.SetOneSlotBusyness(lesson.Slot, true)
	t.LessonChecker.AddLesson(lesson)

	return err
}

func (t *Teacher) CheckLesson(lesson *Lesson) error {
	if err := t.CheckSlot(lesson.Slot); err != nil {
		return err
	}
	if t.IsBusy(lesson.Slot) {
		return fmt.Errorf("teacher is busy")
	}
	if t.CountHourDeficit() <= 0 {
		return fmt.Errorf("enough hours")
	}

	return nil
}
