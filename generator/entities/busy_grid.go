package entities

import (
	"fmt"
)

// NewBusyGrid creates new BusyGrid instance.
//
// It requires a grid with coefficients of comfort (grid).
// Function copies the array before creating a new instance.
func NewBusyGrid(grid [][]float32) *BusyGrid {
	bg := BusyGrid{Grid: make([][]float32, len(grid))}
	for i := range grid {
		bg.Grid[i] = make([]float32, len(grid[i]))
		copy(bg.Grid[i], grid[i])
	}
	return &bg
}

// BusyGrid represents the grid of business for other entities.
type BusyGrid struct {
	Grid [][]float32 // positive - slot is free, negative - slot is busy by lesson, 0 - slot is busy for other reasons
}

// GetFreeSlot returns optimal slot index of the day.
//
// If there are no free slots for both (bg and other) or the lengths are different, returns -1.
// TODO: extract scoring logic and replace index-coupled slice API.
func (bg *BusyGrid) GetOptimalFreeSlot(otherSlots []float32, day int) int {
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

// GetFreeSlots returns filled free slots of the day.
//
// If day isn't within the grid, return an empty array.
// WARNING: slots are represented as float32 values, not as a structure.
func (bg *BusyGrid) GetFreeSlots(day int) (slots []float32) {
	if err := bg.CheckDay(day); err != nil {
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

// MoveLessonTo moves lesson from "from" slot to "to" slot.
// Uses LessonCanBeMoved for check.
func (bg *BusyGrid) MoveLessonTo(from, to LessonSlot) error {
	if err := bg.LessonCanBeMoved(from, to); err != nil {
		return err
	}

	bg.SetSlotBusyState(to, true)
	bg.SetSlotBusyState(from, false)
	return nil
}

// ==========================================================================================================
// ========================================== BUSY STATE MANAGEMENT =========================================
// ==========================================================================================================

// SetSlotBusyState marks a slot as busy or free.
//
// Returns an error if the slot is outside the grid or if it is busy for other reasons.
func (bg *BusyGrid) SetSlotBusyState(slot LessonSlot, isBusy bool) error {
	err := bg.CheckSlot(slot)
	if err != nil {
		return fmt.Errorf("slot is invalid: %s", err.Error())
	}
	if bg.IsBlocked(slot) {
		return fmt.Errorf("can't change slot %s busy state: busy for other reasons", slot.String())
	}

	var sign float32 = -1 // to change sign of coefficient
	if isBusy == bg.IsBusy(slot) {
		sign = 1 // if already done
	}
	bg.Grid[slot.Day][slot.Slot] = sign * bg.Grid[slot.Day][slot.Slot]
	return nil
}

// BlockWeekDay marks all slots of the specified weekday as blocked.
//
// Returns an error if day is not weekday.
func (bg *BusyGrid) BlockWeekDay(day int) error {
	if err := bg.CheckWeekDay(day); err != nil {
		return err
	}

	for week := 0; bg.CheckDay(day+week*7) == nil; week++ {
		err := bg.BlockFullDay(day + week*7)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

// BlockFullDay marks all slots of the day as blocked.
//
// Returns an error if dai isn't within the grid.
func (bg *BusyGrid) BlockFullDay(day int) error {
	if err := bg.CheckDay(day); err != nil {
		return err
	}

	for i := range bg.Grid[day] {
		err := bg.BlockSlot(NewLessonSlot(day, i))
		if err != nil {
			panic(err)
		}
	}

	return nil
}

// BlockSlot marks the slot as blocked.
//
// Returns an error if the slot isn't within the grid.
func (bg *BusyGrid) BlockSlot(slot LessonSlot) error {
	if err := bg.CheckSlot(slot); err != nil {
		return err
	}

	bg.Grid[slot.Day][slot.Slot] = 0.0

	return nil
}

// ==========================================================================================================
// ================================================= CHECKS =================================================
// ==========================================================================================================

// CheckWeekDay checks if the day is a weekday [0-6]. Returns the DayOutError if it is not.
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

// CheckDay checks if the day is within the grid.
// Returns a DayOutError if it is not.
func (bg *BusyGrid) CheckDay(day int) error {
	if len(bg.Grid) <= day || day < 0 {
		return DayOutError{input: day, min: 0, max: len(bg.Grid)}
	}

	return nil
}

// CheckSlot checks if the slot is within a grid.
// Return DayOutError or SlotOutError if it is not.
func (bg *BusyGrid) CheckSlot(slot LessonSlot) error {
	err := bg.CheckDay(slot.Day)
	if err != nil {
		return err
	}

	if len(bg.Grid[slot.Day]) <= slot.Slot || slot.Slot < 0 {
		return SlotOutError{min: 0, max: len(bg.Grid[slot.Day]), input: slot.Slot, day: slot.Day}
	}

	return nil
}

// LessonCanBeMoved checks if the lesson at the "from" slot is marked as busy and the "to" slot is free.
// Returns an error if it is not or if any slot is invalid.
func (bg *BusyGrid) LessonCanBeMoved(from, to LessonSlot) error {
	if err := bg.CheckSlot(from); err != nil {
		return fmt.Errorf("\"from\" slot is invalid: %s", err.Error())
	}
	if !bg.IsBusy(from) {
		return fmt.Errorf("\"from\" slot (%s) is not marked as busy", from.String())
	}

	if err := bg.CheckSlot(to); err != nil {
		return fmt.Errorf("\"to\" slot is invalid: %s", err)
	}

	if bg.IsBusy(to) {
		return fmt.Errorf("\"to\" slot (%s) isn't free", to.String())
	}

	return nil
}

// CheckGapOnAdd checks if the slot is free and adding the lesson does not create a gap.
// Returns an error if it is not.
func (bg *BusyGrid) CheckGapOnAdd(slot LessonSlot) error {
	if !bg.IsFree(slot) {
		return fmt.Errorf("slot (%s) not free", slot.String())
	}

	checkFunc := func(slot LessonSlot, step int) error {
		currentSlot := slot
		currentSlot.Slot += step
		for bg.IsFree(currentSlot) {
			currentSlot.Slot += step
		}
		if bg.CheckSlot(currentSlot) == nil && (currentSlot.Slot-slot.Slot)*step > 1 {
			return fmt.Errorf("slot (%s) creates a window; the nearest not-free slot is %s", slot.String(), currentSlot.String())
		}
		return nil
	}

	if err := checkFunc(slot, 1); err != nil {
		return err
	}
	if err := checkFunc(slot, -1); err != nil {
		return err
	}

	return nil
}

// ==========================================================================================================
// ================================================= STATES =================================================
// ==========================================================================================================

// IsBusy checks if the slot is not free.
// If an error occurs, returns true.
// TODO: use IsFree instead;
func (bg *BusyGrid) IsBusy(slot LessonSlot) bool {
	err := bg.CheckSlot(slot)
	if err != nil {
		return true
	}

	return bg.Grid[slot.Day][slot.Slot] < 0 || bg.IsBlocked(slot)
}

// IsFree checks if the slot is free.
// If an error occurs, returns false.
func (bg *BusyGrid) IsFree(slot LessonSlot) bool {
	if err := bg.CheckSlot(slot); err != nil {
		return false
	}

	return bg.Grid[slot.Day][slot.Slot] > 0
}

// Checks if lesson is at this slot.
// If an error occurs, returns false.
func (bg *BusyGrid) IsLessonOn(slot LessonSlot) bool {
	err := bg.CheckSlot(slot)
	if err != nil {
		return false
	}

	return bg.Grid[slot.Day][slot.Slot] < 0
}

// IsBlocked checks if the slot is blocked.
// If an error occurs, returns true.
func (bg *BusyGrid) IsBlocked(slot LessonSlot) bool {
	if err := bg.CheckSlot(slot); err != nil {
		return true
	}

	return bg.Grid[slot.Day][slot.Slot] == 0.0
}

// ==========================================================================================================
// =============================================== STATISTICS ===============================================
// ==========================================================================================================

// CountWindows returns the sum of windows (gaps between busy slots).
func (bg *BusyGrid) CountWindows() (count int) {
	// Days cycle
	for day := range len(bg.Grid) {
		lastBusy := -1
		// Slots cycle
		for slot := range bg.Grid[day] {
			if !bg.IsFree(NewLessonSlot(day, slot)) {
				if lastBusy != -1 && (slot-lastBusy) > 1 {
					count += slot - lastBusy - 1
				}
				lastBusy = slot
			}
		}
	}
	return
}

// CountLessonsOn returns the sum of lessons on the day.
//
// If day is invalid, returns -1.
func (bg *BusyGrid) CountLessonsOn(day int) (count int) {
	if err := bg.CheckDay(day); err != nil {
		return -1
	}

	for i := range bg.Grid[day] {
		if bg.IsLessonOn(LessonSlot{Day: day, Slot: i}) {
			count++
		}
	}

	return
}

// GetWeekDaysPriority returns slices that contain 7 elements, each representing the priority for the weekdays.
// Priority is calculated as the average slots coefficient on the weekdays.
// WARNING: complex logic.
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

// CountSlotsAtDay returns the sum of free slots on the weekday.
//
// If day isn't a weekday returns 0.
func (bg *BusyGrid) CountSlotsOnDay(day int) (count int) {
	if err := bg.CheckWeekDay(day); err != nil {
		return
	}

	for week := 0; bg.CheckDay(day+week*7) == nil; week++ {
		currentDay := day + week*7
		for slot := range bg.Grid[currentDay] {
			if bg.IsFree(LessonSlot{Day: currentDay, Slot: slot}) {
				count++
			}
		}
	}
	return
}

// CountLessonOverlapping returns the count of overlapping lessons. Counts only lessons that overlap.
func (bg *BusyGrid) CountLessonOverlapping(lessons []*Lesson) (count int) {
	for _, lesson := range lessons {
		// if lesson in not busy slot => overlap or other error
		if !bg.IsLessonOn(lesson.LessonSlot) {
			count++
		}

		// sets slot as free so the next lesson with the same slot wouldn't pass the check
		bg.SetSlotBusyState(lesson.LessonSlot, false)
	}

	// return grid to it first state
	for _, lesson := range lessons {
		bg.SetSlotBusyState(lesson.LessonSlot, true)
	}

	return count
}
