package components

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/generator/services"
)

// MissingLessonsAdder adds missing lessons to the first available day
// in both the teacher's and the student group's schedules.
type MissingLessonsAdder interface {
	GeneratorComponent  // Basic interface for generator component
	AddMissingLessons() // Add a MissingLessonsAdderError to ErrorService
}

// NewMissingLessonAdder creates a MissingLessonsAdder instance.
// It requires an ErrorService, a list of teachers and a LessonService.
func NewMissingLessonAdder(es ErrorService, t []*entities.Teacher, ls services.LessonService) MissingLessonsAdder {
	return &missingLessonsAdder{errorService: es, teachers: t, lessonService: ls}
}

type missingLessonsAdder struct {
	errorService  ErrorService
	teachers      []*entities.Teacher
	lessonService services.LessonService
}

func (ma *missingLessonsAdder) AddMissingLessons() {
	for i := range ma.teachers {
		teacher := ma.teachers[i]

		for _, teacherLoad := range teacher.Load {
			for _, group := range teacherLoad.Groups {
				currentDay := 0
				outOfGrid := false
				for !teacherLoad.Discipline.EnoughHours() && !outOfGrid {
					err := group.CheckDay(currentDay)
					if err != nil {
						outOfGrid = true
						//continue
						break
					}

					for i := range teacher.BusyGrid.Grid[currentDay] {
						slot := entities.LessonSlot{
							Day:  currentDay,
							Slot: i,
						}
						ma.lessonService.AddLesson(
							teacher,
							group,
							teacherLoad.Discipline,
							slot,
							teacherLoad.LessonType,
						)
					}
					delta := group.GetNextDayOfType(teacherLoad.LessonType, currentDay+1)
					if delta == -1 {
						outOfGrid = true
						continue
					}
					currentDay += delta
				}

				if !teacherLoad.Discipline.EnoughHours() {
					ma.errorService.AddError(&MissingLessonsAdderError{
						UnsignedLesson: entities.UnsignedLesson{
							Teacher:      teacher,
							StudentGroup: group,
							Discipline:   teacherLoad.Discipline,
							Type:         teacherLoad.LessonType,
						},
					})
				}
			}
		}
	}
}

// Redirect to AddMissingLessons function
func (ma *missingLessonsAdder) Run() {
	ma.AddMissingLessons()
}

func (ma *missingLessonsAdder) GetErrorService() ErrorService {
	return ma.errorService
}

// MissingLessonsAdderError indicates that the MissingLessonsAdder failed to
// find free slot in the grids for missing lesson.
type MissingLessonsAdderError struct {
	entities.UnsignedLesson
}

func (e *MissingLessonsAdderError) Error() string {
	return fmt.Sprintf("Not enough space of %s or %s for %s %s.",
		e.StudentGroup.Name, e.Teacher.UserName, e.Type.Name, e.Discipline.Name)
}

func (e *MissingLessonsAdderError) GetTypeOfError() GeneratorComponentErrorTypes {
	return MissingLessonsAdderErrorType
}
