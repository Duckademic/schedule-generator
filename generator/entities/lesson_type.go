package entities

import (
	"github.com/google/uuid"
)

type LessonType struct {
	ID    uuid.UUID
	Name  string
	Weeks []int
	Value int
}
