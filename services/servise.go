package services

import (
	"github.com/google/uuid"
)

type Service[T any] interface {
	Create(T) (*T, error)
	Update(T) error
	Find(uuid.UUID) *T
	Delete(uuid.UUID) error
	GetAll() []T
}
