package services

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/types"
)

// StudyLoadService aggregates and manages study loads (UnassignedLessons) that the generator works with.
type StudyLoadService interface {
	GetAll() []*entities.UnassignedLesson // Returns a slice with all study loads as pointers.
}

// NewStudyLoadService creates a new StudyLoadService basic instance.
//
// It requires an array of database study loads (sl), teacher, student group, discipline,
// and lesson type services (ts, sgs, ds, and lts).
//
// Returns an error if any study load is an invalid model.
func NewStudyLoadService(
	sl []types.StudyLoad,
	ts TeacherService,
	sgs StudentGroupService,
	ds DisciplineService,
	lts LessonTypeService,
) (StudyLoadService, error) {
	sls := &studyLoadService{}

	for _, studyLoad := range sl {

		teacher := ts.Find(studyLoad.TeacherID)
		if teacher == nil {
			return nil, fmt.Errorf("teacher %s not found", studyLoad.TeacherID)
		}

		for _, disciplineLoad := range studyLoad.Disciplines {
			discipline := ds.Find(disciplineLoad.DisciplineID)
			if discipline == nil {
				return nil, fmt.Errorf("discipline %s not found", disciplineLoad.DisciplineID)
			}
			lessonType := lts.Find(disciplineLoad.LessonTypeID)
			if lessonType == nil {
				return nil, fmt.Errorf("lesson type %s not found", disciplineLoad.LessonTypeID)
			}

			studentGroups := make([]*entities.StudentGroup, len(disciplineLoad.GroupsID))
			for j, studentGroupID := range disciplineLoad.GroupsID {
				studentGroup := sgs.Find(studentGroupID)
				if studentGroup == nil {
					return nil, fmt.Errorf("student group %s not found", studentGroupID)
				}
				studentGroup.AddLoad(entities.NewStudentLoadKey(discipline, lessonType, teacher), disciplineLoad.Hours)
				for week := range lessonType.Weeks {
					studentGroup.BindWeek(lessonType, week)
				}

				studentGroups[j] = studentGroup
				sls.loads = append(sls.loads, entities.NewUnassignedLesson(lessonType, teacher, studentGroup, discipline))
			}

			// if err := discipline.AddLoad(teacher, disciplineLoad.Hours, studentGroups, lessonType); err != nil {
			// 	return err
			// }
			for _, group := range studentGroups {
				teacher.AddLoad(entities.NewTeacherLoadKey(discipline, group, lessonType), disciplineLoad.Hours)
			}
		}
	}

	return sls, nil
}

// studyLoadService is the basic implementation of the StudyLoadService interface.
type studyLoadService struct {
	loads []*entities.UnassignedLesson
}

func (sls *studyLoadService) GetAll() []*entities.UnassignedLesson {
	return sls.loads
}
