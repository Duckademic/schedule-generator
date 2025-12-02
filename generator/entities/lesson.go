package entities

type LessonSlot struct {
	Day  int
	Slot int
}

// // LessonGroup contains data for correct lesson swapping. Only 1 lesson per week
// type LessonGroup struct {
// 	Structure []*Lesson //index - week
// }

// UnsignedLesson represents a lesson draft without confirmed assignments.
type UnsignedLesson struct {
	Type         *LessonType
	Teacher      *Teacher
	StudentGroup *StudentGroup
	Discipline   *Discipline
}

// Lesson represents an assigned lesson based on an UnsignedLesson.
type Lesson struct {
	UnsignedLesson
	Slot  LessonSlot
	Value int // кількість академічних годин
}

func (l *Lesson) After(other *Lesson) bool {
	if l.Slot.Day > other.Slot.Day {
		return true
	} else if l.Slot.Day < other.Slot.Day {
		return false
	} else if l.Slot.Slot > other.Slot.Slot {
		return true
	}
	return false
}
