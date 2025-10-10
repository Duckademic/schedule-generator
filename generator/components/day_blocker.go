package components

import (
	"fmt"
	"slices"

	"github.com/Duckademic/schedule-generator/generator/entities"
)

// DayBlocker selects days for student groups
type DayBlocker interface {
	SetDayTypes() error // throw an error if at not enough days per group
}

func NewDayBlocker(studentGroups []entities.StudentGroup) DayBlocker {
	db := dayBlocker{}
	db.setGroupExtensions(studentGroups)

	return &db
}

// Extension of group (store data to not calculate every time)
type groupExtension struct {
	group         *entities.StudentGroup // Original StudentGroup
	dayPriorities []float32              // Bigger number - better day for lessons (<0.99 if day is uncomfortable) (length - 7)
	countOfSlots  []int                  // count of slots at every day (length - 7)
	freeDayCount  int                    // count of free (comfortable) days
}

// Time complexity O(1)
func (ge *groupExtension) IsFreeDay(day int) bool {
	return ge.dayPriorities[day] > 0.99
}

func newGroupExtension(group *entities.StudentGroup) *groupExtension {
	ge := groupExtension{
		group:         group,
		dayPriorities: group.GetWeekDaysPriority(),
		countOfSlots:  make([]int, 7),
	}

	for day, _ := range ge.dayPriorities {
		if ge.IsFreeDay(day) {
			ge.freeDayCount++
		}
		ge.countOfSlots[day] = ge.group.CountSlotsAtDay(day)
	}

	return &ge
}

type dayBlocker struct {
	groupExtensions []groupExtension // StudentGroup collection
}

func (db *dayBlocker) SetDayTypes() error {
	daysBlocked := make([]int, 7) // contains num of groups that chose this day

	for _, group := range db.groupExtensions {
		availableDays := []int{0, 1, 2, 3, 4, 5, 6}

		for _, lt := range group.group.GetOwnLessonTypes() {
			//select 2 days for every lesson type
			for tmp_i := 0; tmp_i < 2; tmp_i++ {
				// select day that free and blocked the fewest times
				min := 1000000000
				mIndex := -1
				for _, day := range availableDays {
					if group.IsFreeDay(day) && min > daysBlocked[day] {
						min = daysBlocked[day]
						mIndex = day
					}
				}

				// day not found
				if mIndex == -1 {
					return &SetDayTypeError{
						LessonType:    lt,
						StudentGroup:  group.group,
						DayPriorities: group.dayPriorities,
						AvailableDays: availableDays,
					}
				}

				// if an error ignores this day, deletes it from available days, continues search
				err := group.group.AddDayToLessonType(lt, mIndex)
				if err != nil {
					dayIndex := slices.Index(availableDays, mIndex)
					availableDays = append(availableDays[:dayIndex], availableDays[dayIndex+1:]...)
					tmp_i--
					continue
				}

				// all good, add to blocked day
				daysBlocked[mIndex]++
			}
		}
	}

	return nil
}

func (db *dayBlocker) setGroupExtensions(studentGroups []entities.StudentGroup) {
	db.groupExtensions = make([]groupExtension, len(studentGroups))
	for i := range studentGroups {
		db.groupExtensions[i] = *newGroupExtension(&studentGroups[i])
	}
	// sorts by free day count in increasing order
	slices.SortFunc(db.groupExtensions, func(a, b groupExtension) int {
		if a.freeDayCount == b.freeDayCount {
			return 0
		} else if a.freeDayCount > b.freeDayCount {
			return 1
		}
		return -1
	})
}

type SetDayTypeError struct {
	LessonType    *entities.LessonType
	StudentGroup  *entities.StudentGroup
	DayPriorities []float32
	AvailableDays []int
}

func (e *SetDayTypeError) Error() string {
	return fmt.Sprintf("can't add a day of type %s to group %s", e.LessonType.Name, e.StudentGroup.Name)
}
