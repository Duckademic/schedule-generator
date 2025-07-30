package generator

import "fmt"

func NewBusyGrid(grid [][]bool) *BusyGrid {
	bg := BusyGrid{Grid: make([][]bool, len(grid))}
	for i := range grid {
		bg.Grid[i] = make([]bool, len(grid[i]))
		copy(bg.Grid[i], grid[i])
	}
	return &bg
}

type BusyGrid struct {
	Grid [][]bool
}

func (bg *BusyGrid) SetOneSlotBusyness(slot LessonSlot, isBusy bool) error {
	err := bg.CheckSlot(slot)
	if err != nil {
		return err
	}

	bg.Grid[slot.Day][slot.Slot] = isBusy
	return nil
}

// вільні слоти - то всі незайняті
//
// якщо день за межами сітки (або не матиме слотів) - поверне порожній масив
func (bg *BusyGrid) GetFreeSlots(day int) (slots []bool) {
	err := bg.CheckDay(day)
	if err != nil {
		return
	}

	slots = make([]bool, len(bg.Grid[day]))

	for i := range slots {
		slots[i] = !bg.Grid[day][i]
	}
	return
}

func (bg *BusyGrid) CountWindows() (count int) {
	for i := range len(bg.Grid) {
		lastBusy := -1
		for j, isBusy := range bg.Grid[i] {
			if isBusy {
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

	if len(bg.Grid[slot.Day]) <= slot.Slot {
		return fmt.Errorf("slot %d outside of BusyGrid day %d (max: %d)", slot.Slot, slot.Day, len(bg.Grid[slot.Day]))
	}

	return nil
}

func (bg *BusyGrid) IsBusy(slot LessonSlot) bool {
	err := bg.CheckSlot(slot)
	if err != nil {
		return true
	}

	return bg.Grid[slot.Day][slot.Slot]
}

func (bg *BusyGrid) CountLessonsOn(day int) (count int) {
	for _, isBusy := range bg.Grid[day] {
		if isBusy {
			count++
		}
	}

	return
}
