package components

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/generator/services"
)

// BoneGenerator creates the initial weekly lesson structure (“bone week”)
// by allocating lesson slots for groups and teachers in the first week.
type BoneGenerator interface {
	GeneratorComponent    // Basic interface for generator component
	GenerateBoneLessons() // Add a BoneWeekError to ErrorService if at not enough space at bone week
}

// NewBoneGenerator creates a BoneGenerator instance.
// It requires an ErrorService, a list of study loads, a LessonService.
func NewBoneGenerator(es ErrorService, l []*entities.UnassignedLesson, ls services.LessonService) BoneGenerator {
	return &boneGenerator{errorService: es, loads: l, lessonService: ls}
}

type boneGenerator struct {
	errorService  ErrorService
	loads         []*entities.UnassignedLesson
	lessonService services.LessonService
}

// GenerateBoneLessons allocates lesson slots for the bone week.
// Uses brute force method, starts with teachers, then discipline and student groups,
// then free slots for lesson type.
func (bg *boneGenerator) GenerateBoneLessons() {
	for _, load := range bg.loads {
		teacher := load.Teacher
		studentGroup := load.StudentGroup
		lessonType := load.Type
		// discipline := load.Discipline

		offset := 0
		success := false

		for !success {
			// отримуємо доступний лекційний день
			day := studentGroup.GetNextDayOfType(lessonType, offset)
			if day > 7 || day < 0 {
				// якщо день був не на кістковому тижні, виникає виняток, який треба обробити якось
				bg.errorService.AddError(&BoneWeekError{UnassignedLesson: *load})
				break
			}

			// отримання вільного слота для групи та викладача
			lessonSlot := teacher.GetOptimalFreeSlot(studentGroup.GetFreeSlots(day), day)

			if lessonSlot != -1 {
				slot := entities.LessonSlot{Day: day, Slot: lessonSlot}
				err := bg.lessonService.AssignLesson(*load, slot)
				if err != nil {
					bg.errorService.AddError(NewUnexpectedError("slot is busy but algorithm determined it as free",
						"boneGenerator", "GenerateBoneLessons", &FalseFreeSlotError{
							UnassignedLesson: *load,
							slot:             slot,
							err:              err,
						}))
				}
				success = true
			}
			offset = day + 1
		}

	}
}

// Redirect to GenerateBoneLessons function
func (bg *boneGenerator) Run() {
	bg.GenerateBoneLessons()
}

func (bg *boneGenerator) GetErrorService() ErrorService {
	return bg.errorService
}

// BoneWeekError indicates that the BoneGenerator failed to allocate
// enough space for lessons within the bone week.
type BoneWeekError struct {
	entities.UnassignedLesson
}

func (e *BoneWeekError) Error() string {
	return fmt.Sprintf("Not enough space in bone week of %s or %s for %s %s.",
		e.StudentGroup.Name, e.Teacher.UserName, e.Type.Name, e.Discipline.Name)
}

func (e *BoneWeekError) GetTypeOfError() GeneratorComponentErrorTypes {
	return BoneWeekErrorType
}

// FalseFreeSlotError indicates that slot is busy but algorithm determined it as free.
type FalseFreeSlotError struct {
	entities.UnassignedLesson
	slot entities.LessonSlot
	err  error
}

func (e *FalseFreeSlotError) Error() string {
	return fmt.Sprintf("false free slot %d/%d of %s or %s grid for %s %s. error: %s", e.slot.Day, e.slot.Slot,
		e.StudentGroup.Name, e.Teacher.UserName, e.Type.Name, e.Discipline.Name, e.err.Error())
}
