package generator

import (
	"slices"

	"github.com/Duckademic/schedule-generator/types"
	"github.com/google/uuid"
)

type StudentGroup struct {
	ID          uuid.UUID
	Name        string
	BusyGrid    [][]bool
	LectureDays []int
}

type StudentGroupService interface {
	GetAll() []StudentGroup
	SetOneSlotBusyness(*StudentGroup, LessonSlot, bool) error
	GetFreeSlots(group *StudentGroup, day int) []bool
	GetLectureDay(group *StudentGroup, startDay int) int // ПЕРЕРОБИТИ НА УЗАГАЛЬНЕННЯ
	CountLessonsOn(group *StudentGroup, day int) int
	Find(uuid.UUID) *StudentGroup
}

type studentGroupService struct {
	studentGroups    []StudentGroup
	maxLessonsPerDay int
}

func NewStudentGroupService(studentGroups []types.StudentGroup, maxLessonsPerDay int, busyGrid [][]bool) (StudentGroupService, error) {
	sgs := studentGroupService{
		studentGroups:    make([]StudentGroup, len(studentGroups)),
		maxLessonsPerDay: maxLessonsPerDay,
	}

	for i := range studentGroups {
		sgs.studentGroups[i] = StudentGroup{ID: studentGroups[i].ID, Name: studentGroups[i].Name}
		sgs.studentGroups[i].BusyGrid = make([][]bool, len(busyGrid))
		for j := range busyGrid {
			sgs.studentGroups[i].BusyGrid[j] = make([]bool, len(busyGrid[j]))
			copy(sgs.studentGroups[i].BusyGrid[j], busyGrid[j])
		}
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

	sgs.maxLessonsPerDay = maxLessonsPerDay

	return &sgs, nil
}

func (sgs *studentGroupService) GetAll() []StudentGroup {
	return sgs.studentGroups
}

func (sgs *studentGroupService) SetOneSlotBusyness(group *StudentGroup, slot LessonSlot, isBusy bool) error {
	group.BusyGrid[slot.Day][slot.Slot] = isBusy
	return nil
}

func (sgs *studentGroupService) GetFreeSlots(group *StudentGroup, day int) (slots []bool) {
	slots = make([]bool, len(group.BusyGrid[day]))

	// випадок, коли ще немає занять
	if sgs.CountLessonsOn(group, day) == 0 {
		for i := range slots {
			slots[i] = true
		}
		return
	}

	for i := range group.BusyGrid[day] {
		// пропускаємо 1 елемент щоб далі не виникло помилок
		if i == 0 {
			continue
		}

		// якщо у поточному слоті вже є пара, а у попередньому ні, вписуємо попередній слот як доступний
		if group.BusyGrid[day][i] {
			if !group.BusyGrid[day][i-1] {
				slots[i-1] = true
			}
			// якщо у слоті немає пари, а у попередньому вона є, то вписуємо поточний слот як доступний
		} else {
			if group.BusyGrid[day][i-1] {
				slots[i] = true
			}
		}
	}
	return
}

// returns -1 if student group hasn't free lecture day
func (sgs *studentGroupService) GetLectureDay(group *StudentGroup, startDay int) int {
	if group == nil {
		return -1
	}

	for i := startDay; i < len(group.BusyGrid); i++ {
		if slices.Contains(group.LectureDays, i%7) {
			if sgs.CountLessonsOn(group, i) < sgs.maxLessonsPerDay {
				return i
			}
		}
	}

	return -1
}

func (sgs *studentGroupService) CountLessonsOn(group *StudentGroup, day int) (count int) {
	if group == nil {
		return 0
	}

	for _, isBusy := range group.BusyGrid[day] {
		if isBusy {
			count++
		}
	}

	return
}

// return will be nil if not found
func (sgs *studentGroupService) Find(id uuid.UUID) *StudentGroup {
	for i := range sgs.studentGroups {
		if sgs.studentGroups[i].ID == id {
			return &sgs.studentGroups[i]
		}
	}

	return nil
}
