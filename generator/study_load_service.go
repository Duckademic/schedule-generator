package generator

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/types"
)

type StudyLoad struct {
	Teacher     *Teacher
	Disciplines []DisciplineLoad
}

type DisciplineLoad struct {
	Discipline *Discipline
	Groups     []*StudentGroup
	Hours      int
}

type StudyLoadService interface {
	GetAll() []StudyLoad
}

func NewStudyLoadService(studyLoads []types.StudyLoad, ts TeacherService, sgs StudentGroupService) (StudyLoadService, error) {
	sls := studyLoadService{studyLoads: make([]StudyLoad, len(studyLoads))}

	for i, studyLoad := range studyLoads {
		sls.studyLoads[i] = StudyLoad{
			Teacher: ts.Find(studyLoad.TeacherID),
		}
		if sls.studyLoads[i].Teacher == nil {
			return nil, fmt.Errorf("teacher %s not found", studyLoad.TeacherID)
		}

		for _, disciplineLoad := range studyLoad.Disciplines {
			dl := DisciplineLoad{
				Discipline: &Discipline{ID: disciplineLoad.DisciplineID},
				Hours:      disciplineLoad.Hours,
				Groups:     make([]*StudentGroup, len(disciplineLoad.GroupsID)),
			}

			for j, studentGroupID := range disciplineLoad.GroupsID {
				dl.Groups[j] = sgs.Find(studentGroupID)
				if dl.Groups[j] == nil {
					return nil, fmt.Errorf("student group %s not found", studentGroupID)
				}
			}

			sls.studyLoads[i].Disciplines = append(sls.studyLoads[i].Disciplines, dl)
		}
	}

	sls.currentStudyLoad = &sls.studyLoads[0]

	return &sls, nil
}

type studyLoadService struct {
	studyLoads       []StudyLoad
	currentStudyLoad *StudyLoad
}

func (sls *studyLoadService) GetAll() []StudyLoad {
	return sls.studyLoads
}
