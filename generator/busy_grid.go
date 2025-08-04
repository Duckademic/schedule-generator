package generator

import (
	"fmt"
)

func NewBusyGrid(grid [][]float32) *BusyGrid {
	bg := BusyGrid{Grid: make([][]float32, len(grid))}
	for i := range grid {
		bg.Grid[i] = make([]float32, len(grid[i]))
		copy(bg.Grid[i], grid[i])
	}
	return &bg
}

type BusyGrid struct {
	Grid [][]float32 // додатнє - вільне, від'ємне - зайняте заняттям, 0 - зайняте (інші причини)
}

func (bg *BusyGrid) SetOneSlotBusyness(slot LessonSlot, isBusy bool) error {
	err := bg.CheckSlot(slot)
	if err != nil {
		return err
	}

	var sign float32 = -1
	if isBusy == bg.IsBusy(slot) {
		sign = 1
	}
	bg.Grid[slot.Day][slot.Slot] = sign * bg.Grid[slot.Day][slot.Slot]
	return nil
}

// вільні слоти - то всі незайняті
//
// якщо день за межами сітки (або не матиме слотів) - поверне порожній масив
func (bg *BusyGrid) GetFreeSlots(day int) (slots []float32) {
	err := bg.CheckDay(day)
	if err != nil {
		return
	}

	slots = make([]float32, len(bg.Grid[day]))

	for i := range slots {
		if !bg.IsBusy(LessonSlot{Day: day, Slot: i}) {
			slots[i] = bg.Grid[day][i]
		}
	}
	return
}

func (bg *BusyGrid) CountWindows() (count int) {
	for i := range len(bg.Grid) {
		lastBusy := -1
		for j := range bg.Grid[i] {
			if bg.IsBusy(LessonSlot{Day: i, Slot: j}) {
				if lastBusy != -1 && (j-lastBusy) > 1 {
					count += j - lastBusy - 1
				}
				lastBusy = j
			}
		}
	}
	return
}

type DayOutError struct {
	min   int
	max   int
	input int
}

func (d DayOutError) Error() string {
	return fmt.Sprintf("day %d outside of BusyGrid (%d to %d)", d.input, d.min, d.max)
}

func (bg *BusyGrid) CheckDay(day int) error {
	if len(bg.Grid) <= day || day < 0 {
		return DayOutError{input: day, min: 0, max: len(bg.Grid)}
	}

	return nil
}

func (bg *BusyGrid) CheckSlot(slot LessonSlot) error {
	err := bg.CheckDay(slot.Day)
	if err != nil {
		return err
	}

	if len(bg.Grid[slot.Day]) <= slot.Slot || slot.Slot < 0 {
		return fmt.Errorf("slot %d outside of BusyGrid day %d (max: %d)", slot.Slot, slot.Day, len(bg.Grid[slot.Day]))
	}

	return nil
}

func (bg *BusyGrid) IsBusy(slot LessonSlot) bool {
	err := bg.CheckSlot(slot)
	if err != nil {
		return true
	}

	return bg.Grid[slot.Day][slot.Slot] <= 0
}

func (bg *BusyGrid) CountLessonsOn(day int) (count int) {
	for i := range bg.Grid[day] {
		if bg.IsBusy(LessonSlot{Day: day, Slot: i}) {
			count++
		}
	}

	return
}

// returns -1 if there are not free slot for both (bg and other) or length are different
func (bg *BusyGrid) GetFreeSlot(otherSlots []float32, day int) int {
	if err := bg.CheckDay(day); err != nil {
		return -1
	}
	if len(bg.Grid[day]) != len(otherSlots) {
		return -1
	}

	var max float32 = 0.0
	maxI := -1
	for i := range bg.Grid[day] {
		if !bg.IsBusy(LessonSlot{Day: day, Slot: i}) {
			value := bg.Grid[day][i] * otherSlots[i]
			if max < value {
				maxI = i
				max = value
			}
		}
	}

	return maxI
}

// returns slices which contains 7 elements
func (bg *BusyGrid) GetWeekDaysPriority() (result []float32) {
	result = make([]float32, 7)
	for day := range 7 {
		for week := 0; bg.CheckDay(day+week*7) == nil; week++ {
			currentDay := day + week*7
			var average float32 = 0
			for slot, value := range bg.Grid[currentDay] {
				average = ((average * float32(slot)) + value) / (float32(slot) + 1)
			}

			result[day] = (result[day]*float32(week) + average) / (float32(week) + 1)
		}
	}
	return
}

func (bg *BusyGrid) CountSlotsAtDay(day int) (count int) {
	if err := bg.CheckWeekDay(day); err != nil {
		return
	}

	for week := 0; bg.CheckDay(day+week*7) == nil; week++ {
		currentDay := day + week*7
		for slot := range bg.Grid[currentDay] {
			if !bg.IsBusy(LessonSlot{Day: currentDay, Slot: slot}) {
				count++
			}
		}
	}
	return
}

func (bg *BusyGrid) SetDayBusyness(newBusyness []float32, day int) error {
	if err := bg.CheckWeekDay(day); err != nil {
		return err
	}
	if len(newBusyness) != len(bg.Grid[day]) {
		return fmt.Errorf("incorrect length of the new busyness (%d instead of %d)", len(newBusyness), len(bg.Grid[day]))
	}

	for week := 0; bg.CheckDay(day+week*7) == nil; week++ {
		copy(bg.Grid[day], newBusyness)
	}

	return nil
}

func (bg *BusyGrid) CheckWeekDay(day int) error {
	if day < 0 || day > 6 {
		return DayOutError{
			min:   0,
			max:   6,
			input: day,
		}
	}

	return nil
}

func (bg *BusyGrid) CountLessonOverlapping(lessons []*Lesson) (count int) {
	for _, lesson := range lessons {
		if bg.Grid[lesson.Slot.Day][lesson.Slot.Slot] >= 0 {
			count++
		}

		bg.Grid[lesson.Slot.Day][lesson.Slot.Slot] = -bg.Grid[lesson.Slot.Day][lesson.Slot.Slot]
	}

	for _, lesson := range lessons {
		bg.Grid[lesson.Slot.Day][lesson.Slot.Slot] = -bg.Grid[lesson.Slot.Day][lesson.Slot.Slot]
	}

	return count
}
