package generator

import (
	"fmt"
	"slices"

	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type StudentGroup struct {
	BusyGrid
	ID               uuid.UUID
	Name             string
	MaxLessonsPerDay int
	DaysOfType       map[*LessonType][]int
}

func (sg *StudentGroup) IsBusy(slot LessonSlot) bool {
	if err := sg.CheckSlot(slot); err != nil {
		return true
	}

	return sg.CountLessonsOn(slot.Day) >= sg.MaxLessonsPerDay || !sg.GetFreeSlots(slot.Day)[slot.Slot]
}

func (sg *StudentGroup) GetFreeSlots(day int) (slots []bool) {
	if err := sg.CheckDay(day); err != nil {
		return []bool{}
	}

	slots = make([]bool, len(sg.Grid[day]))

	// випадок, коли ще немає занять
	if sg.CountLessonsOn(day) == 0 {
		for i := range slots {
			slots[i] = true
		}
		return
	}

	for i := range sg.Grid[day] {
		// пропускаємо 1 елемент щоб далі не виникло помилок
		if i == 0 {
			continue
		}

		// якщо у поточному слоті вже є пара, а у попередньому ні, вписуємо попередній слот як доступний
		if sg.Grid[day][i] {
			if !sg.Grid[day][i-1] {
				slots[i-1] = true
			}
			// якщо у слоті немає пари, а у попередньому вона є, то вписуємо поточний слот як доступний
		} else {
			if sg.Grid[day][i-1] {
				slots[i] = true
			}
		}
	}
	return
}

// returns -1 if student group hasn't free day
func (sg *StudentGroup) GetNextDayOfType(lType *LessonType, startDay int) int {
	if len(sg.DaysOfType[lType]) == 0 {
		return -1
	}

	for i := startDay; i < len(sg.Grid); i++ {
		if slices.Contains(sg.DaysOfType[lType], i%7) {
			if sg.CountLessonsOn(i) < sg.MaxLessonsPerDay {
				return i
			}
		}
	}

	return -1
}

func (sg *StudentGroup) SetDayType(lType *LessonType, day int) error {
	if day < 0 || day > 6 {
		return fmt.Errorf("day %d out of range (%d to %d)", day, 0, 6)
	}

	days := sg.DaysOfType[lType]
	if slices.Contains(days, day) {
		return fmt.Errorf("day %d already typed as %s", day, lType.Name)
	}

	sg.DaysOfType[lType] = append(days, day)
	slices.Sort(sg.DaysOfType[lType])
	return nil
}

type StudentGroupService interface {
	GetAll() []StudentGroup
	Find(uuid.UUID) *StudentGroup
	CountWindows() int
}

type studentGroupService struct {
	studentGroups []StudentGroup
}

func NewStudentGroupService(studentGroups []types.StudentGroup, maxLessonsPerDay int, busyGrid [][]bool) (StudentGroupService, error) {
	sgs := studentGroupService{
		studentGroups: make([]StudentGroup, len(studentGroups)),
	}

	for i := range studentGroups {
		sgs.studentGroups[i] = StudentGroup{
			ID:               studentGroups[i].ID,
			Name:             studentGroups[i].Name,
			MaxLessonsPerDay: maxLessonsPerDay,
			DaysOfType:       map[*LessonType][]int{},
		}
		sgs.studentGroups[i].BusyGrid = *NewBusyGrid(busyGrid)
	}

	return &sgs, nil
}

func (sgs *studentGroupService) GetAll() []StudentGroup {
	return sgs.studentGroups
}

// return nil if not found
func (sgs *studentGroupService) Find(id uuid.UUID) *StudentGroup {
	for i := range sgs.studentGroups {
		if sgs.studentGroups[i].ID == id {
			return &sgs.studentGroups[i]
		}
	}

	return nil
}

func (sgs *studentGroupService) CountWindows() (count int) {
	for _, g := range sgs.studentGroups {
		count += g.CountWindows()
	}
	return
}
