package entities

import "fmt"

// SlotOutError is returned when a slot in specific day is out of range.
//
// Error: slot %input% at %day% day outside of BusyGrid (%min% to %max%)
type SlotOutError struct {
	min   int
	max   int
	input int
	day   int
}

func (s SlotOutError) Error() string {
	return fmt.Sprintf("slot %d at %d day outside of BusyGrid (%d to %d)", s.input, s.day, s.min, s.max)
}

// DayOutError is returned when a day value is out of range.
//
// Error: day %input% outside of BusyGrid (%min% to %max%)
type DayOutError struct {
	min   int
	max   int
	input int
}

func (d DayOutError) Error() string {
	return fmt.Sprintf("day %d outside of BusyGrid (%d to %d)", d.input, d.min, d.max)
}
