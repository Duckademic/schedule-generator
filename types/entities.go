package types

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Model interface {
	FindID() string
}

func FindFirst(arr []Model, model Model) int {
	for i, m := range arr {
		if m.FindID() == model.FindID() {
			return i
		}
	}
	return -1
}

type StudyLoad struct {
	Teacher     Teacher
	Disciplines []DisciplineLoad
}

func (sl *StudyLoad) FindID() string {
	panic("not implemented")
}

type DisciplineLoad struct {
	Discipline Discipline
	Groups     []StudentGroup
	Hours      int
}

// ==============================================================

type StudentGroup struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	// Number string // номер групи (32)

	Business    [][]bool `json:"-"`
	LectureDays []int    `json:"-"`
}

func (sg *StudentGroup) FindID() string {
	return sg.Name
}

type Teacher struct {
	ID       uuid.UUID `json:"id"`
	UserName string    `json:"user_name"`
	// AcademicDegree string // асистент/доцент/професор

	Business [][]bool `json:"-"`
}

func (t *Teacher) FindID() string {
	return t.UserName
}

type Discipline struct {
	ID   uuid.UUID
	Name string
	// Lessons map[string]int // тип - кількість годин
}

func (d *Discipline) FindID() string {
	return d.Name
}

type Lesson struct {
	ID        uuid.UUID
	StartTime time.Time
	EndTime   time.Time
	Value     int // кількість академічних годин
	Type      LessonType
	// Gap       int
}

func CheckWindows(teachers []Teacher, groups []StudentGroup) (teacherW, groupW int) {
	for i := range len(teachers[0].Business) {
		for _, t := range teachers {
			lastBusy := -1
			for j, isBusy := range t.Business[i] {
				if isBusy {
					if lastBusy != -1 && (j-lastBusy) > 1 {
						teacherW += j - lastBusy - 1
					}
					lastBusy = j
				}
			}
		}

		for _, t := range groups {
			lastBusy := -1
			for j, isBusy := range t.Business[i] {
				if isBusy {
					if lastBusy != -1 && (j-lastBusy) > 1 {
						groupW += j - lastBusy - 1
					}
					lastBusy = j
				}
			}
		}
	}
	return
}

func TestCheckWindows() {
	// Умова:
	// День 0: заняття у 0-й та 2-й парах → 1 "вікно"
	// День 1: заняття у 1-й та 4-й → 2 "вікна" (між 1 і 4: 2 і 3)
	teacher := Teacher{
		Business: [][]bool{
			{true, false, true, false},        // День 0
			{false, true, false, false, true}, // День 1
		},
	}
	group := StudentGroup{
		Business: [][]bool{
			{true, false, false, true},       // День 0 → 2 "вікна"
			{true, false, true, false, true}, // День 1 → 2 "вікна"
		},
	}

	teachers := []Teacher{teacher}
	groups := []StudentGroup{group}

	wantTeacherW := 3 // 1 (день 0) + 2 (день 1)
	wantGroupW := 4   // 2 (день 0) + 2 (день 1)

	gotTeacherW, gotGroupW := CheckWindows(teachers, groups)

	if gotTeacherW != wantTeacherW || gotGroupW != wantGroupW {
		panic(fmt.Sprintf("Expected (teacherW=%d, groupW=%d), got (teacherW=%d, groupW=%d)",
			wantTeacherW, wantGroupW, gotTeacherW, gotGroupW))
	}
	panic("all correct")
}

type LessonType struct {
	Name string
}
