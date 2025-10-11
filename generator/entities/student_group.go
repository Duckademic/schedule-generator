package entities

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
)

type StudentGroupLoad struct {
	Days  []int
	Weeks []int
	LessonChecker
}

type StudentGroup struct {
	BusyGrid
	ID                uuid.UUID
	Name              string
	MaxLessonsPerDay  int
	LessonTypeBinding map[*LessonType]*StudentGroupLoad
	Teachers          []*Teacher
}

// ==========================================================================================================
// ========================================== BusyGrid OVERRIDES ============================================
// ==========================================================================================================

func (sg *StudentGroup) IsBusy(slot LessonSlot) bool {
	if err := sg.BusyGrid.CheckSlot(slot); err != nil {
		return true
	}

	slotIsBusy := true
	// якщо поточний слот вільний, то один з сусідніх має бути зайнятим, причому в сітці, або всі вільні
	if !sg.BusyGrid.IsBusy(LessonSlot{Day: slot.Day, Slot: slot.Slot}) {
		if sg.CountLessonsOn(slot.Day) == 0 {
			slotIsBusy = false
		} else {
			for _, value := range []int{-1, 1} {
				tmpSlot := LessonSlot{Day: slot.Day, Slot: slot.Slot + value}
				if err := sg.CheckSlot(tmpSlot); err == nil && sg.BusyGrid.IsBusy(tmpSlot) {
					slotIsBusy = false
				}
			}
		}
	}

	return sg.CountLessonsOn(slot.Day) >= sg.MaxLessonsPerDay || slotIsBusy
}

func (sg *StudentGroup) CountSlotsAtDay(day int) (count int) {
	if day < 0 || day > 6 {
		return
	}

	for week := 0; sg.CheckDay(day+week*7) == nil; week++ {
		currentDay := day + week*7
		delta := 0
		for slot := range sg.Grid[currentDay] {
			if !sg.IsBusy(LessonSlot{Day: currentDay, Slot: slot}) {
				delta++
			}
		}
		if delta > sg.MaxLessonsPerDay {
			delta = sg.MaxLessonsPerDay
		}
		count += delta
	}
	return
}

func (sg *StudentGroup) GetFreeSlots(day int) (slots []float32) {
	if err := sg.CheckDay(day); err != nil {
		return []float32{}
	}

	slots = make([]float32, len(sg.Grid[day]))

	// випадок, коли ще немає занять
	if sg.CountLessonsOn(day) == 0 {
		for i := range slots {
			slots[i] = sg.Grid[day][i]
		}
		return
	}

	for i := range sg.Grid[day] {
		// пропускаємо 1 елемент щоб далі не виникло помилок
		if i == 0 {
			continue
		}

		// якщо у поточному слоті вже є пара, а у попередньому ні, вписуємо попередній слот як доступний
		if sg.IsBusy(LessonSlot{Day: day, Slot: i}) {
			if !sg.IsBusy(LessonSlot{Day: day, Slot: i - 1}) {
				slots[i-1] = sg.Grid[day][i-1]
			}
			// якщо у слоті немає пари, а у попередньому вона є, то вписуємо поточний слот як доступний
		} else {
			if sg.IsBusy(LessonSlot{Day: day, Slot: i - 1}) {
				slots[i] = sg.Grid[day][i]
			}
		}
	}
	return
}

func (sg *StudentGroup) MoveLessonTo(lesson *Lesson, to LessonSlot) error {
	if err := sg.LessonCanBeMoved(lesson, to); err != nil {
		return err
	}

	return sg.BusyGrid.MoveLessonTo(lesson, to)
}

// Additionally checks the type of day
func (sg *StudentGroup) LessonCanBeMoved(lesson *Lesson, to LessonSlot) error {
	if err := sg.BusyGrid.LessonCanBeMoved(lesson, to); err != nil {
		return err
	}

	if !sg.IsDayOfType(lesson.Type, to.Day) {
		return fmt.Errorf("%d is not day of the type %s", to.Day, lesson.Type.Name)
	}

	return nil
}

// ==========================================================================================================
// ========================================== LessonType BINDING ============================================
// ==========================================================================================================

// returns -1 if student group hasn't free day
func (sg *StudentGroup) GetNextDayOfType(lType *LessonType, startDay int) int {
	if len(sg.LessonTypeBinding[lType].Days) == 0 {
		return -1
	}

	for i := startDay; i < len(sg.Grid); i++ {
		if sg.IsDayOfType(lType, i) {
			if sg.CountLessonsOn(i) < sg.MaxLessonsPerDay {
				return i
			}
		}
	}

	return -1
}

func (sg *StudentGroup) IsDayOfType(lType *LessonType, day int) bool {
	for lessonType, load := range sg.LessonTypeBinding {
		if lessonType != lType && slices.Contains(load.Weeks, day/7) {
			return false
		}
	}

	return slices.Contains(sg.LessonTypeBinding[lType].Days, day%7) || slices.Contains(sg.LessonTypeBinding[lType].Weeks, day/7)
}

func (sg *StudentGroup) GetMaxHours(lType *LessonType) int {
	if sgl, ok := sg.LessonTypeBinding[lType]; ok {
		return sgl.RequiredHours
	}
	return 0
}

func (sg *StudentGroup) AddBindingToLessonType(lType *LessonType, hours int, teacher *Teacher) error {
	if lType == nil {
		return fmt.Errorf("lesson type is nil")
	}

	_, ok := sg.LessonTypeBinding[lType]
	if !ok {
		sg.LessonTypeBinding[lType] = &StudentGroupLoad{}
	}

	sg.LessonTypeBinding[lType].RequiredHours += hours
	if slices.Index(sg.Teachers, teacher) == -1 {
		sg.Teachers = append(sg.Teachers, teacher)
	}
	return nil
}

func (sg *StudentGroup) AddDayToLessonType(lType *LessonType, day int) error {
	if day < 0 || day > 6 {
		return fmt.Errorf("day %d out of range (%d to %d)", day, 0, 6)
	}

	load, ok := sg.LessonTypeBinding[lType]
	if !ok {
		return fmt.Errorf("type %s not found", lType.Name)
	}
	for lessonType, load := range sg.LessonTypeBinding {
		if slices.Contains(load.Days, day) {
			return fmt.Errorf("day %d already typed as %s", day, lessonType.Name)
		}
	}

	load.Days = append(load.Days, day)
	slices.Sort(load.Days)
	return nil
}

func (sg *StudentGroup) AddWeekToLessonType(lType *LessonType, week int) error {
	if len(sg.Grid)/7 < week || week < 0 {
		return fmt.Errorf("week %d out of range (%d to %d)", week, 0, len(sg.Grid)/7)
	}

	load, ok := sg.LessonTypeBinding[lType]
	if !ok {
		return fmt.Errorf("type %s not found", lType.Name)
	}
	for lessonType, load := range sg.LessonTypeBinding {
		if slices.Contains(load.Weeks, week) {
			return fmt.Errorf("week %d already typed as %s", week, lessonType.Name)
		}
	}

	load.Weeks = append(load.Weeks, week)
	slices.Sort(load.Weeks)
	return nil
}

func (sg *StudentGroup) AddLesson(lesson *Lesson, ignoreCheck bool) error {
	err := sg.CheckLesson(lesson)
	if err != nil && !ignoreCheck {
		return err
	}

	sg.SetOneSlotBusyness(lesson.Slot, true)
	sg.LessonTypeBinding[lesson.Type].AddLesson(lesson)

	return err
}

func (sg *StudentGroup) CheckLesson(lesson *Lesson) error {
	if err := sg.CheckSlot(lesson.Slot); err != nil {
		return err
	}
	if sg.IsBusy(lesson.Slot) {
		return fmt.Errorf("student group is busy")
	}

	if !sg.IsDayOfType(lesson.Type, lesson.Slot.Day) {
		return fmt.Errorf("type %s not in the correct day", lesson.Type.Name)
	}

	return nil
}

func (sg *StudentGroup) CountHourDeficit() (count int) {
	for _, studentGroupLoad := range sg.LessonTypeBinding {
		count += studentGroupLoad.CountHourDeficit()
	}

	return count
}

// Returns sum of lesson overlap.
func (sg *StudentGroup) CountLessonOverlapping() (count int) {
	for _, load := range sg.LessonTypeBinding {
		count += sg.BusyGrid.CountLessonOverlapping(load.Lessons)
	}

	return
}

// Returns types of lesson for student group
func (sg *StudentGroup) GetOwnLessonTypes() []*LessonType {
	keys := make([]*LessonType, 0, len(sg.LessonTypeBinding))
	for lt := range sg.LessonTypeBinding {
		keys = append(keys, lt)
	}
	return keys
}
