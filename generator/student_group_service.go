package generator

import (
	"fmt"
	"slices"

	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type StudentGroupLoad struct {
	Days []int
	LessonChecker
}

type StudentGroup struct {
	BusyGrid
	ID               uuid.UUID
	Name             string
	MaxLessonsPerDay int
	DaysOfType       map[*LessonType]*StudentGroupLoad
}

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

// returns -1 if student group hasn't free day
func (sg *StudentGroup) GetNextDayOfType(lType *LessonType, startDay int) int {
	if len(sg.DaysOfType[lType].Days) == 0 {
		return -1
	}

	for i := startDay; i < len(sg.Grid); i++ {
		if slices.Contains(sg.DaysOfType[lType].Days, i%7) {
			if sg.CountLessonsOn(i) < sg.MaxLessonsPerDay {
				return i
			}
		}
	}

	return -1
}

func (sg *StudentGroup) GetMaxHours(lType *LessonType) int {
	if sgl, ok := sg.DaysOfType[lType]; ok {
		return sgl.RequiredHours
	}
	return 0
}

func (sg *StudentGroup) AddDayType(lType *LessonType, hours int) error {
	if lType == nil {
		return fmt.Errorf("lesson type is nil")
	}

	_, ok := sg.DaysOfType[lType]
	if !ok {
		sg.DaysOfType[lType] = &StudentGroupLoad{}
	}

	sg.DaysOfType[lType].RequiredHours += hours
	return nil
}

func (sg *StudentGroup) SetDayType(lType *LessonType, day int) error {
	if day < 0 || day > 6 {
		return fmt.Errorf("day %d out of range (%d to %d)", day, 0, 6)
	}

	load, ok := sg.DaysOfType[lType]
	if !ok {
		return fmt.Errorf("type %s not found", lType.Name)
	}
	if slices.Contains(load.Days, day) {
		return fmt.Errorf("day %d already typed as %s", day, lType.Name)
	}

	load.Days = append(load.Days, day)
	slices.Sort(load.Days)
	return nil
}

func (sg *StudentGroup) AddLesson(lesson *Lesson, ignoreCheck bool) error {
	err := sg.CheckLesson(lesson)
	if err != nil && !ignoreCheck {
		return err
	}

	sg.SetOneSlotBusyness(lesson.Slot, true)
	sg.DaysOfType[lesson.Type].AddLesson(lesson)

	return err
}

func (sg *StudentGroup) CheckLesson(lesson *Lesson) error {
	if err := sg.CheckSlot(lesson.Slot); err != nil {
		return err
	}
	if sg.IsBusy(lesson.Slot) {
		return fmt.Errorf("student group is busy")
	}
	if sg.DaysOfType[lesson.Type].CountHourDeficit() <= 0 {
		return fmt.Errorf("enough hours of type %s", lesson.Type.Name)
	}

	return nil
}

func (sg *StudentGroup) CountHourDeficit() (count int) {
	for _, studentGroupLoad := range sg.DaysOfType {
		count += studentGroupLoad.CountHourDeficit()
	}

	return count
}

func (sg *StudentGroup) CountLessonOverlapping() (count int) {
	for _, load := range sg.DaysOfType {
		count += sg.BusyGrid.CountLessonOverlapping(load.Lessons)
	}

	return
}

type StudentGroupService interface {
	GetAll() []StudentGroup
	Find(uuid.UUID) *StudentGroup
	CountWindows() int
	CountHourDeficit() int
	CountLessonOverlapping() int
}

type studentGroupService struct {
	studentGroups []StudentGroup
}

func NewStudentGroupService(studentGroups []types.StudentGroup, maxLessonsPerDay int, busyGrid [][]float32) (StudentGroupService, error) {
	sgs := studentGroupService{
		studentGroups: make([]StudentGroup, len(studentGroups)),
	}

	for i := range studentGroups {
		sgs.studentGroups[i] = StudentGroup{
			ID:               studentGroups[i].ID,
			Name:             studentGroups[i].Name,
			MaxLessonsPerDay: maxLessonsPerDay,
			DaysOfType:       map[*LessonType]*StudentGroupLoad{},
		}
		studentGroup := &sgs.studentGroups[i]
		studentGroup.BusyGrid = *NewBusyGrid(busyGrid)

		md := studentGroups[i].MilitaryDay - 1
		if md != -1 {
			if err := studentGroup.CheckWeekDay(md); err != nil {
				return nil, err
			}
			studentGroup.SetDayBusyness(make([]float32, len(studentGroup.Grid[md])), md)
		}
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

func (sgs *studentGroupService) CountHourDeficit() (count int) {
	for _, studentGroup := range sgs.studentGroups {
		count += studentGroup.CountHourDeficit()
	}

	return count
}

func (sgs *studentGroupService) CountLessonOverlapping() (count int) {
	for _, studentGroup := range sgs.studentGroups {
		count += studentGroup.CountLessonOverlapping()
	}

	return
}
