package components

import (
	"fmt"
	"slices"

	"github.com/Duckademic/schedule-generator/generator/entities"
)

// DayBlocker selects days for student groups
type DayBlocker interface {
	GeneratorComponent // Basic interface for generator component
	SetDayTypes()      // Add a SetDayTypeError to ErrorService if at not enough days per group
}

// NewDayBlocker creates a DayBlocker instance.
// It requires an ErrorService and a list of student groups.
func NewDayBlocker(studentGroups []*entities.StudentGroup, errorService ErrorService) DayBlocker {
	db := dayBlocker{errorService: errorService}
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

func (ge *groupExtension) IsFreeDay(day int) bool {
	return ge.dayPriorities[day] > 0.99
}

func newGroupExtension(group *entities.StudentGroup) *groupExtension {
	ge := groupExtension{
		group:         group,
		dayPriorities: group.GetWeekDaysPriority(),
		countOfSlots:  make([]int, 7),
	}

	for day := range ge.dayPriorities {
		if ge.IsFreeDay(day) {
			ge.freeDayCount++
		}
		ge.countOfSlots[day] = ge.group.CountSlotsAtDay(day)
	}

	return &ge
}

type dayBlocker struct {
	groupExtensions []groupExtension // StudentGroup collection
	errorService    ErrorService     // Collection for errors
}

func (db *dayBlocker) SetDayTypes() {
	daysBlocked := make([]int, 7) // contains num of groups that chose this day

	for _, group := range db.groupExtensions {
		availableDays := []int{0, 1, 2, 3, 4, 5, 6}

		for _, lt := range group.group.GetOwnLessonTypes() {
			//select 2 days for every lesson type
			for tmp_i := 0; tmp_i < 2; tmp_i++ { // break after error assigned
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
					db.errorService.AddError(&SetDayTypeError{
						LessonType:    lt,
						StudentGroup:  group.group,
						DayPriorities: group.dayPriorities,
						AvailableDays: availableDays,
					})
					break // continue with next group
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
}

func (db *dayBlocker) GetErrorService() ErrorService {
	return db.errorService
}

// Redirect to SetDayTypes function
func (db *dayBlocker) Run() {
	db.SetDayTypes()
}

func (db *dayBlocker) setGroupExtensions(studentGroups []*entities.StudentGroup) {
	db.groupExtensions = make([]groupExtension, len(studentGroups))
	for i := range studentGroups {
		db.groupExtensions[i] = *newGroupExtension(studentGroups[i])
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

func (e *SetDayTypeError) GetTypeOfError() GeneratorComponentErrorTypes {
	return SetDayTypeErrorType
}
