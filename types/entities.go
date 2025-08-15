package types

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID `json:"id" binding:"required" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}

type StudyLoad struct {
	TeacherID   uuid.UUID
	Disciplines []DisciplineLoad
}

type DisciplineLoad struct {
	DisciplineID uuid.UUID
	GroupsID     []uuid.UUID
	Hours        int
	LessonTypeID uuid.UUID
}

// ==============================================================

type StudentGroup struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required,min=4"`
	MilitaryDay int       `json:"military_day" binding:"gte=1,lte=7"`
	// Number string // номер групи (32)
}

type Teacher struct {
	Model
	UserName string `json:"user_name" binding:"required,min=4,max=64" gorm:"type:varchar(64);unique"`
	Priority int    `json:"priority"`
	// масив з бажаннями викладача
	// AcademicDegree string // асистент/доцент/професор
}

type Discipline struct {
	ID   uuid.UUID
	Name string
	// Lessons map[string]int // тип - кількість годин
}

type Lesson struct {
	ID        uuid.UUID  `json:"id" validate:"required"`
	StartTime time.Time  `json:"start_time" binding:"required"`
	EndTime   time.Time  `json:"end_time" binding:"required"`
	Value     int        `json:"value" binding:"required,gt=0"` // кількість академічних годин
	Type      LessonType `json:"type" binding:"required"`
	// Gap       int
}

type LessonType struct {
	ID    uuid.UUID `json:"id" binding:"required"`
	Name  string    `json:"name" binding:"required,min=4"`
	Weeks int       `json:"weeks"` // кількість тижнів на початку навчання заповнених тільки цими типами занять
}
