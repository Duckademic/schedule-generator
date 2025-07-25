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
	return sg.CountLessonsOn(slot.Day) >= sg.MaxLessonsPerDay || sg.BusyGrid.IsBusy(slot)
}

// returns -1 if student group hasn't free lecture day
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

// func (sgs *studentGroupService) GetFreeSlots(group *StudentGroup, day int) (slots []bool) {
// 	slots = make([]bool, len(group.BusyGrid[day]))

// 	// випадок, коли ще немає занять
// 	if sgs.CountLessonsOn(group, day) == 0 {
// 		for i := range slots {
// 			slots[i] = true
// 		}
// 		return
// 	}

// 	for i := range group.BusyGrid[day] {
// 		// пропускаємо 1 елемент щоб далі не виникло помилок
// 		if i == 0 {
// 			continue
// 		}

// 		// якщо у поточному слоті вже є пара, а у попередньому ні, вписуємо попередній слот як доступний
// 		if group.BusyGrid[day][i] {
// 			if !group.BusyGrid[day][i-1] {
// 				slots[i-1] = true
// 			}
// 			// якщо у слоті немає пари, а у попередньому вона є, то вписуємо поточний слот як доступний
// 		} else {
// 			if group.BusyGrid[day][i-1] {
// 				slots[i] = true
// 			}
// 		}
// 	}
// 	return
// }

// return will be nil if not found
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
