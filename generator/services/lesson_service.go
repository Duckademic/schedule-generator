package services

import (
	"fmt"
	"sort"

	"github.com/Duckademic/schedule-generator/generator/entities"
)

// LessonService aggregates and manages lessons that the generator works with.
type LessonService interface {
	GetAll() []*entities.Lesson // Returns a slice with all lessons as pointers.
	// Assigns a lesson to the selected slot.
	AssignLesson(entities.UnassignedLesson, entities.LessonSlot) error
	MoveLessonTo(*entities.Lesson, entities.LessonSlot) error // MoveLessonTo moves lesson to another slot (to).
	GetWeekLessons(int) []*entities.Lesson                    // TODO: collect bone lessons in another structure.
}

// NewLessonService creates a new LessonService basic instance.
//
// It requires a number of academic hours for lessons (lesson value - lv).
//
// Returns an error if the lesson value is below or equal to zero.
func NewLessonService(lv int) (LessonService, error) {
	if lv <= 0 {
		return nil, fmt.Errorf("lessonValue below/equal to 0 (%d)", lv)
	}

	ls := lessonService{lessonValue: lv}

	return &ls, nil
}

// lessonService is the basic implementation of the LessonService interface.
type lessonService struct {
	lessons     []*entities.Lesson
	lessonValue int
}

func (ls *lessonService) GetAll() []*entities.Lesson {
	return ls.lessons
}
func (ls *lessonService) AssignLesson(ul entities.UnassignedLesson, slot entities.LessonSlot) error {
	if err := ul.Validate(); err != nil {
		return err
	}

	lesson := entities.NewLesson(ul, slot, ls.lessonValue)

	if err := ul.Teacher.CheckLesson(lesson); err != nil {
		return err
	}
	if err := ul.StudentGroup.CheckLesson(lesson); err != nil {
		return err
	}

	ls.lessons = append(ls.lessons, lesson)

	if err := lesson.StudentGroup.AddLesson(lesson); err != nil {
		panic("pass the check before, but error accurse")
	}
	if err := lesson.Teacher.AddLesson(lesson); err != nil {
		panic("pass the check before, but error accurse")
	}

	return nil
}
func (ls *lessonService) GetWeekLessons(week int) (res []*entities.Lesson) {
	for _, l := range ls.lessons {
		if l.Day/7 == week {
			res = append(res, l)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].Day != res[j].Day {
			return res[i].Day < res[j].Day
		}
		return res[i].Slot < res[j].Slot
	})
	return
}
func (ls *lessonService) MoveLessonTo(lesson *entities.Lesson, to entities.LessonSlot) error {
	if err := lesson.Teacher.LessonCanBeMoved(lesson.LessonSlot, to); err != nil {
		return err
	}
	if err := lesson.StudentGroup.LessonCanBeMoved(lesson, to); err != nil {
		return err
	}

	if err := lesson.Teacher.MoveLessonTo(lesson.LessonSlot, to); err != nil {
		panic("pass the check before, but error accurse")
	}
	if err := lesson.StudentGroup.MoveLessonTo(lesson, to); err != nil {
		panic("pass the check before, but error accurse")
	}
	lesson.MoveLessonTo(to)
	return nil
}
