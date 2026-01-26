package entities

import (
	"fmt"

	"github.com/google/uuid"
)

// LessonType represents type of lessons.
type LessonType struct {
	ID          uuid.UUID // Unique identifier of the LessonType.
	Name        string    // Human-readable identifier of the LessonType.
	Weeks       []int     // List of week numbers when only this type can be assigned.
	Value       int       // Number of academic hours assigned to this LessonType
	DayRequired int       // <-- NOT RESPONSIBILITY OF THIS OBJECT
}

// LessonTypeBinder stores bindings between lesson types and calendar.
type LessonTypeBinder interface {
	// Assigns a lesson type to a specific week.
	//
	// Returns an error if the week is already blocked.
	BindWeek(*LessonType, int) error
	// Assigns a lesson type to a specific weekday.
	//
	// Returns an error if the weekday is already blocked.
	BindWeekday(*LessonType, int) error
	// Checks whether the given day matches the lesson type.
	//
	// Week binding has higher priority than weekday binding.
	IsDayOfType(*LessonType, int) bool
}

// NewLessonTypeBinder creates a new basic LessonTypeChecker instance.
func NewLessonTypeBinder() LessonTypeBinder {
	return &lessonTypeBinder{
		weekBinding: make(map[int]*LessonType),
		dayBinding:  make([]*LessonType, 7),
	}
}

// lessonTypeBinder is the basic implementation of the LessonTypeBlocker interface.
type lessonTypeBinder struct {
	weekBinding map[int]*LessonType
	dayBinding  []*LessonType
}

func (c *lessonTypeBinder) BindWeek(lt *LessonType, week int) error {
	if lt, ok := c.weekBinding[week]; ok {
		return fmt.Errorf("week %d already blocked (by %s)", week, lt.Name)
	}

	c.weekBinding[week] = lt

	return nil
}
func (c *lessonTypeBinder) BindWeekday(lt *LessonType, day int) error {
	if !c.IsWeekday(day) {
		return fmt.Errorf("day %d is not the number of the weekday", day)
	}

	if c.dayBinding[day] != nil {
		return fmt.Errorf("day %d already blocked (by %s)", day, c.dayBinding[day].Name)
	}

	c.dayBinding[day] = lt

	return nil
}
func (c *lessonTypeBinder) IsDayOfType(lt *LessonType, day int) bool {
	trueLT, ok := c.weekBinding[day/7]
	if ok {
		return trueLT == lt
	}

	return c.dayBinding[day%7] == lt
}
func (c *lessonTypeBinder) IsWeekday(day int) bool {
	return day >= 0 && day <= 6
}
