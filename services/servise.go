package services

import (
	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type Service[T any] interface {
	Create(T) (*T, error)
	Update(T) error
	Find(uuid.UUID) *T
	Delete(uuid.UUID) error
	GetAll() []T
}

type SimpleService[T types.Model] struct {
	objects []T
}

func (s *SimpleService[T]) Create(t T) error {
	s.objects = append(s.objects, t)

	return nil
}

func (s *SimpleService[T]) ReadFirst(t T) T {
	for _, o := range s.objects {
		if o.FindID() == t.FindID() {
			return o
		}
	}
	panic(t.FindID() + " not found")
}

func (s *SimpleService[T]) GetAll() []T {
	return s.objects
}
