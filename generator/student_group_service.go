package generator

import (
	"slices"

	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type StudentGroup struct {
	BusyGrid
	ID               uuid.UUID
	Name             string
	LectureDays      []int
	MaxLessonsPerDay int
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

// ПЕРЕРОБИТИ returns -1 if student group hasn't free lecture day
func (sg *StudentGroup) GetLectureDay(startDay int) int {
	for i := startDay; i < len(sg.BusyGrid.Grid); i++ {
		if slices.Contains(sg.LectureDays, i%7) {
			if sg.CountLessonsOn(i) < sg.MaxLessonsPerDay {
				return i
			}
		}
	}

	return -1
}

// ПЕРЕРОБИТИ/ВИНЕСТИ в BusyGrid
func (sg *StudentGroup) CountLessonsOn(day int) (count int) {
	for _, isBusy := range sg.BusyGrid.Grid[day] {
		if isBusy {
			count++
		}
	}

	return
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
		}
		sgs.studentGroups[i].BusyGrid = *NewBusyGrid(busyGrid)
	}

	// if len(sgs.studentGroups) >= 6 {
	// 	sgs.studentGroups[0].LectureDays = []int{1, 2}
	// 	sgs.studentGroups[1].LectureDays = []int{2, 3}
	// 	sgs.studentGroups[2].LectureDays = []int{3, 4}
	// 	sgs.studentGroups[3].LectureDays = []int{4, 5}
	// 	sgs.studentGroups[4].LectureDays = []int{5, 1}
	// 	sgs.studentGroups[5].LectureDays = []int{1, 2}
	// }
	for i := range studentGroups {
		sgs.studentGroups[i].LectureDays = []int{1, 2}
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
