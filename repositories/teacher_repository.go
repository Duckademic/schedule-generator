package repositories

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/types"
	"gorm.io/gorm"
)

type TeacherRepository interface {
	Repository[types.Teacher]
}

func NewTeacherRepository(db *gorm.DB) (TeacherRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("database is nil")
	}

	tr := teacherRepository{
		simpleRepository: simpleRepository[types.Teacher]{
			db: db,
		},
	}

	if err := tr.Migrate(); err != nil {
		return nil, fmt.Errorf("teacher model migration error: %s", err)
	}

	return &tr, nil
}

type teacherRepository struct {
	simpleRepository[types.Teacher]
}
