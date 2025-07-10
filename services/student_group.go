package services

import (
	"slices"

	"github.com/Duckademic/schedule-generator/types"
)

type StudentGroupServise struct {
	studentGroups    []types.StudentGroup
	currentGroup     *types.StudentGroup
	maxLessonsPerDay int
}

func NewStudentGroupService(studentGroups []types.StudentGroup, maxLessonsPerDay int) (*StudentGroupServise, error) {
	sgs := StudentGroupServise{studentGroups: studentGroups}
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
	sgs.currentGroup = &sgs.studentGroups[0]

	sgs.maxLessonsPerDay = maxLessonsPerDay

	return &sgs, nil
}

func (sgs *StudentGroupServise) GetAll() []types.StudentGroup {
	return sgs.studentGroups
}

func (sgs *StudentGroupServise) SetBusyness(free [][]bool) {
	for i := range sgs.studentGroups {
		sgs.studentGroups[i].Business = make([][]bool, len(free))
		for j := range free {
			sgs.studentGroups[i].Business[j] = make([]bool, len(free[j]))
			copy(sgs.studentGroups[i].Business[j], free[j])
		}
	}
}

// НЕПРОТЕСТОВАНА ====================================================================
func (sgs *StudentGroupServise) SetOneSlotBusyness(groupId string, day, slot int, isBusy bool) {
	group := sgs.Find(groupId)
	group.Business[day][slot] = isBusy
}

// НЕПРОТЕСТОВАНА ====================================================================
func (sgs *StudentGroupServise) GetFreeSlots(groupId string, day int) (slots []bool) {
	group := sgs.Find(groupId)
	slots = make([]bool, len(group.Business[day]))

	// випадок, коли ще немає занять
	if sgs.CountLessonsOn(groupId, day) == 0 {
		for i := range slots {
			slots[i] = true
		}
		return
	}

	for i := range group.Business[day] {
		// пропускаємо 1 елемент щоб далі не виникло помилок
		if i == 0 {
			continue
		}

		// якщо у поточному слоті вже є пара, а у попередньому ні, вписуємо попередній слот як доступний
		if group.Business[day][i] {
			if !group.Business[day][i-1] {
				slots[i-1] = true
			}
			// якщо у слоті немає пари, а у попередньому вона є, то вписуємо поточний слот як доступний
		} else {
			if group.Business[day][i-1] {
				slots[i] = true
			}
		}
	}
	return
}

// returns -1 if student group hasn't free lecture day
func (sgs *StudentGroupServise) GetLectureDay(groupId string, startDay int) int {
	group := sgs.Find(groupId)
	if group == nil {
		panic("student group nor found")
	}

	for i := startDay; i < len(group.Business); i++ {
		if slices.Contains(group.LectureDays, i%7) {
			if sgs.CountLessonsOn(groupId, i) < sgs.maxLessonsPerDay {
				return i
			}
		}
	}

	return -1
}

func (sgs *StudentGroupServise) CountLessonsOn(groupId string, day int) (count int) {
	group := sgs.Find(groupId)

	for _, isBusy := range group.Business[day] {
		if isBusy {
			count++
		}
	}

	return
}

// return will be nil if not found
func (sgs *StudentGroupServise) Find(id string) *types.StudentGroup {
	var group *types.StudentGroup
	if sgs.currentGroup.Name != id {
		for i := range sgs.studentGroups {
			if sgs.studentGroups[i].Name == id {
				group = &sgs.studentGroups[i]
				break
			}
		}
		sgs.currentGroup = group
	} else {
		group = sgs.currentGroup
	}

	return group
}
