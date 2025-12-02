package components

import (
	"fmt"

	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/generator/services"
)

// BoneGenerator creates the initial weekly lesson structure (“bone week”)
// by allocating lesson slots for groups and teachers.
type BoneGenerator interface {
	GeneratorComponent    // Basic interface for generator component
	GenerateBoneLessons() // Add a BoneWeekError to ErrorService if at not enough space at bone week
}

// NewBoneGenerator creates a BoneGenerator instance.
// It requires an ErrorService, a list of teachers, a LessonService,
// and the target bone-week number (bw).
func NewBoneGenerator(es ErrorService, t []*entities.Teacher, ls services.LessonService, bw int) BoneGenerator {
	return &boneGenerator{errorService: es, teachers: t, lessonService: ls, boneWeek: bw}
}

type boneGenerator struct {
	errorService  ErrorService
	teachers      []*entities.Teacher
	lessonService services.LessonService
	boneWeek      int
}

// GenerateBoneLessons allocates lesson slots for the bone week.
// Uses brute force method, starts with teachers, then discipline and student groups,
// then free slots for lesson type.
func (bg *boneGenerator) GenerateBoneLessons() {
	for i := range bg.teachers {
		teacher := bg.teachers[i]

		for _, teacherLoad := range teacher.Load {
			for _, studentGroup := range teacherLoad.Groups {
				offset := 0
				success := false

				for !success {
					// отримуємо доступний лекційний день
					day := studentGroup.GetNextDayOfType(teacherLoad.LessonType, bg.boneWeek*7+offset)
					if day > bg.boneWeek*7+7 || day < 0 {
						// якщо день був не на кістковому тижні, виникає виняток, який треба обробити якось
						bg.errorService.AddError(
							&BoneWeekError{
								teacher:    teacher,
								group:      studentGroup,
								lessonType: teacherLoad.LessonType,
								discipline: teacherLoad.Discipline,
							},
						)
					}

					// отримання вільного слота для групи та викладача
					lessonSlot := teacher.GetFreeSlot(studentGroup.GetFreeSlots(day), day)

					if lessonSlot != -1 {
						slot := entities.LessonSlot{Day: day, Slot: lessonSlot}
						err := bg.lessonService.AddLesson(teacher, studentGroup, teacherLoad.Discipline, slot, teacherLoad.LessonType)
						if err != nil {
							bg.errorService.AddError(NewUnexpectedError("slot is busy but algorithm determined it as free",
								"boneGenerator", "GenerateBoneLessons", &FalseFreeSlotError{
									UnsignedLesson: entities.UnsignedLesson{
										Teacher:      teacher,
										StudentGroup: studentGroup,
										Discipline:   teacherLoad.Discipline,
										Type:         teacherLoad.LessonType,
									},
									slot: slot,
									err:  err,
								}))
						}
						success = true
					}
					offset = day - bg.boneWeek*7 + 1
				}

			}
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
	teacher    *entities.Teacher
	group      *entities.StudentGroup
	lessonType *entities.LessonType
	discipline *entities.Discipline
}

func (e *BoneWeekError) Error() string {
	return fmt.Sprintf("Not enough space in bone week of %s or %s for %s %s.",
		e.group.Name, e.teacher.UserName, e.lessonType.Name, e.discipline.Name)
}

func (e *BoneWeekError) GetTypeOfError() GeneratorComponentErrorTypes {
	return BoneWeekErrorType
}

// FalseFreeSlotError indicates that slot is busy but algorithm determined it as free.
type FalseFreeSlotError struct {
	entities.UnsignedLesson
	slot entities.LessonSlot
	err  error
}

func (e *FalseFreeSlotError) Error() string {
	return fmt.Sprintf("false free slot %d/%d of %s or %s grid for %s %s. error: %s", e.slot.Day, e.slot.Slot,
		e.StudentGroup.Name, e.Teacher.UserName, e.Type.Name, e.Discipline.Name, e.err.Error())
}
