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
}

// looks for free slots to selected lessons. move lesson to it if found
func (imp *improver) ImproveToNext() bool {
	lessons := imp.lessonService.GetAll()
	// runs until finds free slot or be out of lessons
	for {
		currentLesson := lessons[imp.currentLesson]
		dayOutOfRange := false
		startSlot := currentLesson.Slot.Slot
		for day := currentLesson.Slot.Day; !dayOutOfRange; day++ {
			slotOutOfRange := false
			for slot := startSlot; !slotOutOfRange && !dayOutOfRange; slot++ {
				err := imp.lessonService.MoveLesson(currentLesson, entities.LessonSlot{Slot: slot, Day: day})
				switch err.(type) {
				case entities.DayOutError:
					dayOutOfRange = true
				case entities.SlotOutError:
					slotOutOfRange = true
				case nil:
					return true
				}
			}
			startSlot = 0
		}

		if currentLesson.Slot != imp.startSlot {
			imp.lessonService.MoveLesson(currentLesson, imp.startSlot)
		}

		imp.currentLesson++
		if imp.currentLesson >= len(lessons) {
			return false
		}
		imp.startSlot = lessons[imp.currentLesson].Slot
	}
}

func (imp *improver) SubmitChanges() {
	lessons := imp.lessonService.GetAll()
	imp.startSlot = lessons[imp.currentLesson].Slot
}
