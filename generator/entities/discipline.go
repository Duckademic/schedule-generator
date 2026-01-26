package entities

import "github.com/google/uuid"

// Discipline represents a university subject in the scheduling context.
type Discipline struct {
	ID   uuid.UUID // Unique identifier of the Discipline.
	Name string    // Human-readable identifier of the Discipline.
}

// NewDiscipline creates a new Discipline instance.
//
// It requires discipline's id and name.
func NewDiscipline(id uuid.UUID, name string) *Discipline {
	return &Discipline{
		ID:   id,
		Name: name,
	}
}
