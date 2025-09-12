package repositories

import (
	"fmt"
	"os"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(*T) error
	GetFirst(uuid.UUID) *T
	Update(*T) error
	Delete(*T) error
	GetAll() []T
}

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_USER_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %s", err.Error())
	}

	return db, nil
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to close database: %s", err.Error())
	}

	err = sqlDB.Close()
	if err != nil {
		return fmt.Errorf("failed to close database: %s", err.Error())
	}

	return nil
}

type simpleRepository[T any] struct {
	db *gorm.DB
}

func (sr *simpleRepository[T]) Create(obj *T) error {
	sr.db.Create(obj)
	return nil
}

func (sr *simpleRepository[T]) GetFirst(id uuid.UUID) *T {
	var obj T
	sr.db.First(&obj, "id = ?", id)

	var zero T
	if reflect.DeepEqual(obj, zero) {
		return nil
	}
	return &obj
}

func (sr *simpleRepository[T]) Update(obj *T) error {
	sr.db.Save(obj)
	return nil
}

func (sr *simpleRepository[T]) Delete(obj *T) error {
	sr.db.Delete(obj)
	return nil
}

func (sr *simpleRepository[T]) GetAll() (objects []T) {
	sr.db.Set("gorm:auto_preload", true).Find(&objects)
	return
}

func (sr *simpleRepository[T]) Migrate() error {
	var zero T
	return sr.db.AutoMigrate(&zero)
}
