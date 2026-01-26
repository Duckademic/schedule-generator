package components

import (
	"github.com/Duckademic/schedule-generator/generator/entities"
	"github.com/Duckademic/schedule-generator/generator/services"
)

// Improver improves finished schedule
type Improver interface {
	ImproveToNext() bool // Improve improves schedule. Returns false if there are not available improvements
	SubmitChanges()      // SubmitChanges submits changes after previous submit
}

func NewImprover(lessonService services.LessonService) Improver {
	return &improver{lessonService: lessonService}
}

type improver struct {
	lessonService services.LessonService
	currentLesson int
	startSlot     entities.LessonSlot // home slot for current lesson
	currentSlot   entities.LessonSlot // start slot for improving current lesson
}

// looks for free slots to selected lessons. move lesson to it if found
func (imp *improver) ImproveToNext() bool {
	lessons := imp.lessonService.GetAll()
	// runs until finds free slot or be out of lessons
	for {
		imp.currentSlot.Slot += 1 // moves to the next slot instead of keeping the lesson in the same one
		currentLesson := lessons[imp.currentLesson]
		dayOutOfRange := false
		startSlot := imp.currentSlot.Slot // for the first entry, this value should match the current slot value
		for day := imp.currentSlot.Day; !dayOutOfRange; day++ {
			slotOutOfRange := false
			for slot := startSlot; !slotOutOfRange && !dayOutOfRange; slot++ {
				err := imp.lessonService.MoveLessonTo(currentLesson, entities.LessonSlot{Slot: slot, Day: day})
				switch err.(type) {
				case entities.DayOutError:
					dayOutOfRange = true
				case entities.SlotOutError:
					slotOutOfRange = true
				case nil:
					imp.currentSlot = currentLesson.LessonSlot
					return true
				}
			}
			startSlot = 0 // starts slots from beginning
		}

		if currentLesson.LessonSlot != imp.startSlot {
			imp.lessonService.MoveLessonTo(currentLesson, imp.startSlot)
		}

		imp.currentLesson++
		if imp.currentLesson >= len(lessons) {
			return false
		}
		imp.startSlot = lessons[imp.currentLesson].LessonSlot
		imp.currentSlot = entities.LessonSlot{Day: 0, Slot: 0}
	}
}

func (imp *improver) SubmitChanges() {
	lessons := imp.lessonService.GetAll()
	imp.startSlot = lessons[imp.currentLesson].LessonSlot
}
