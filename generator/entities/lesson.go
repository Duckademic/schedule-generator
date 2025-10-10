package entities

type LessonSlot struct {
	Day  int
	Slot int
}

type Lesson struct {
	Slot         LessonSlot
	Value        int // кількість академічних годин
	Type         *LessonType
	Teacher      *Teacher
	StudentGroup *StudentGroup
	Discipline   *Discipline
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
