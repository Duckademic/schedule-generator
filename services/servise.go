package services

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/repositories"
	"github.com/google/uuid"
)

type Service[T comparable] interface {
	Create(T) (*T, error)
	Update(T) error
	Find(uuid.UUID) *T
	Delete(uuid.UUID) error
	GetAll() []T
}

type gormSimpleService[T comparable] struct {
	repo repositories.Repository[T]
}

func (ss *gormSimpleService[T]) Create(obj T) (*T, error) {
	return &obj, ss.repo.Create(&obj)
}

func (ss *gormSimpleService[T]) Update(obj T) error {
	return ss.repo.Update(&obj)
}

func (ss *gormSimpleService[T]) Find(id uuid.UUID) *T {
	return ss.repo.GetFirst(id)
}

func (ss *gormSimpleService[T]) Delete(id uuid.UUID) error {
	obj := ss.Find(id)
	if obj == nil {
		return fmt.Errorf("can't delete, object %s not found", id)
	}

	return ss.repo.Delete(obj)
}

func (ss *gormSimpleService[T]) GetAll() []T {
	return ss.repo.GetAll()
}
