package types

import (
	"time"

	"github.com/google/uuid"
)

type StudyLoad struct {
	TeacherID   uuid.UUID
	Disciplines []DisciplineLoad
}

type DisciplineLoad struct {
	DisciplineID uuid.UUID
	GroupsID     []uuid.UUID
	Hours        int
}

// ==============================================================

type StudentGroup struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	// Number string // номер групи (32)
}

type Teacher struct {
	ID       uuid.UUID `json:"id"`
	UserName string    `json:"user_name"`
	// AcademicDegree string // асистент/доцент/професор
}

type Discipline struct {
	ID   uuid.UUID
	Name string
	// Lessons map[string]int // тип - кількість годин
}

type Lesson struct {
	ID        uuid.UUID  `json:"id"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
	Value     int        `json:"value"` // кількість академічних годин
	Type      LessonType `json:"type"`
	// Gap       int
}

type LessonType struct {
	Name string `json:"name"`
}
