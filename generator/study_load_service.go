package generator

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/types"
)

func LoadStudyLoads(
	studyLoads []types.StudyLoad,
	ts TeacherService,
	sgs StudentGroupService,
	ds DisciplineService,
	lts LessonTypeService,
) error {
	for _, studyLoad := range studyLoads {
		teacher := ts.Find(studyLoad.TeacherID)
		if teacher == nil {
			return fmt.Errorf("teacher %s not found", studyLoad.TeacherID)
		}

		for _, disciplineLoad := range studyLoad.Disciplines {
			discipline := ds.Find(disciplineLoad.DisciplineID)
			if discipline == nil {
				return fmt.Errorf("discipline %s not found", disciplineLoad.DisciplineID)
			}
			lessonType := lts.Find(disciplineLoad.LessonTypeID)
			if lessonType == nil {
				return fmt.Errorf("lesson type %s not found", disciplineLoad.LessonTypeID)
			}

			dl := DisciplineLoad{
				Teacher:    teacher,
				LoadHours:  len(disciplineLoad.GroupsID) * disciplineLoad.Hours,
				Groups:     make([]*StudentGroup, len(disciplineLoad.GroupsID)),
				LessonType: lts.Find(disciplineLoad.LessonTypeID),
			}

			tl := TeacherLoad{
				Discipline: discipline,
				LessonType: lessonType,
				Groups:     dl.Groups,
			}

			for j, studentGroupID := range disciplineLoad.GroupsID {
				dl.Groups[j] = sgs.Find(studentGroupID)
				if dl.Groups[j] == nil {
					return fmt.Errorf("student group %s not found", studentGroupID)
				}
			}

			if err := discipline.AddLoad(&dl); err != nil {
				return err
			}
			if err := teacher.AddLoad(&tl); err != nil {
				return err
			}
		}
	}

	return nil
}
