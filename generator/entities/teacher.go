package entities

import (
	"fmt"

	"github.com/google/uuid"
)

// Teacher represents a university teacher in the scheduling context.
//
// The model enforces teaching load constraints for groups and disallows
// simultaneous classes.
//
// TODO: add teacher availability constraints.
type Teacher struct {
	BusyGrid                     // Availability grid.
	TeacherLoadService           // Handles teacher load validation logic.
	ID                 uuid.UUID // Unique identifier of the Teacher.
	UserName           string    // Human-readable identifier of the Teacher.
	Priority           int       // Higher value means higher priority (used for sorting).
}

// NewTeacher creates a new Teacher instance.
//
// It requires teacher's id, name (un), priority (p), busy grid for teacher (bg), and load service (tls).
func NewTeacher(id uuid.UUID, un string, p int, bg *BusyGrid, tls TeacherLoadService) *Teacher {
	return &Teacher{
		BusyGrid:           *bg,
		TeacherLoadService: tls,
		ID:                 id,
		UserName:           un,
		Priority:           p,
	}
}

// NewTeacher creates a new Teacher instance with default configuration.
//
// It requires teacher's id, name (un), priority (p) and busy grid for teacher (bg).
func NewDefaultTeacher(id uuid.UUID, un string, p int, bg *BusyGrid) *Teacher {
	return NewTeacher(id, un, p, bg, NewTeacherLoadService())
}

// AddLesson register the lesson.
//
// Uses CheckLesson for check.
func (t *Teacher) AddLesson(lesson *Lesson) error {
	err := t.CheckLesson(lesson)
	if err != nil {
		return err
	}

	t.SetSlotBusyState(lesson.LessonSlot, true)
	t.TeacherLoadService.AddLesson(lesson)

	return err
}

// CheckLesson checks if the lesson can be added. It checks slot validation, availability and load limits.
//
// Return an error if validation fails.
func (t *Teacher) CheckLesson(lesson *Lesson) error {
	if err := t.CheckSlot(lesson.LessonSlot); err != nil {
		return err
	}
	if !t.IsFree(lesson.LessonSlot) {
		return fmt.Errorf("teacher is busy")
	}
	if t.IsEnoughLessons() {
		return fmt.Errorf("teacher %s has enough hours", t.UserName)
	}

	return nil
}

// TeacherLoadService tracks and evaluates the study workload for Teacher.
type TeacherLoadService interface {
	LoadService                            // Basic interface for load validation logic.
	AddLoad(key TeacherLoadKey, hours int) // Registers a new required load entry.
	// Returns true if the teacher doesn't require additional lessons for the specific load.
	IsEnoughLessonsFor(TeacherLoadKey) bool
}

// NewTeacherLoadService creates a new TeacherLoadService basic instance.
func NewTeacherLoadService() TeacherLoadService {
	return &teacherLoadService{
		loads: make(map[TeacherLoadKey]teacherLoad),
	}
}

// teacherLoadService is the basic implementation of the TeacherLoadService interface.
type teacherLoadService struct {
	loads map[TeacherLoadKey]teacherLoad
}

func (s *teacherLoadService) AddLesson(lesson *Lesson) {
	key := TeacherLoadKey{
		discipline:   lesson.Discipline,
		studentGroup: lesson.StudentGroup,
		lessonType:   lesson.Type,
	}

	load, ok := s.loads[key]
	if ok {
		load.checker.AddLesson(lesson)
	} else {
		panic("load not found")
	}
}
func (s *teacherLoadService) CountHourDeficit() (count int) {
	for _, load := range s.loads {
		count += load.checker.CountHourDeficit()
	}

	return
}
func (s *teacherLoadService) IsEnoughLessons() bool {
	for _, load := range s.loads {
		if !load.checker.IsEnoughLessons() {
			return false
		}
	}

	return true
}
func (s *teacherLoadService) GetAssignedLessons() (result []*Lesson) {
	for _, load := range s.loads {
		result = append(result, load.checker.GetAssignedLessons()...)
	}

	return
}
func (s *teacherLoadService) AddLoad(key TeacherLoadKey, hours int) {
	_, ok := s.loads[key]
	if !ok {
		s.loads[key] = teacherLoad{checker: NewLoadService(hours)}
	}
}
func (s *teacherLoadService) IsEnoughLessonsFor(key TeacherLoadKey) bool {
	load, ok := s.loads[key]
	if !ok {
		return true
	}

	return load.checker.IsEnoughLessons()
}

// NewTeacherLoadKey creates a new TeacherLoadKey instance.
//
// It requires pointers to discipline, student group, and lesson type.
func NewTeacherLoadKey(d *Discipline, sg *StudentGroup, lt *LessonType) TeacherLoadKey {
	return TeacherLoadKey{discipline: d, studentGroup: sg, lessonType: lt}
}

// TeacherLoadKey is a composite key used to identify a teacher load entry.
type TeacherLoadKey struct {
	discipline   *Discipline
	studentGroup *StudentGroup
	lessonType   *LessonType
}

type teacherLoad struct {
	checker LoadService
}
