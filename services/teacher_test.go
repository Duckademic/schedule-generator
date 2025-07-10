package services

import (
	"reflect"
	"testing"

	"github.com/Duckademic/schedule-generator/types"
)

func TestNewTeacherServiceAndGetAll(t *testing.T) {
	teachers := []types.Teacher{
		{
			UserName: "Ivan Ivanov",
		},
	}

	ts, err := NewTeacherService(teachers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := ts.GetAll()

	if !reflect.DeepEqual(got, teachers) {
		t.Errorf("expected %v, got %v", teachers, got)
	}
}

func TestTeacherSetBusyness(t *testing.T) {
	teachers := []types.Teacher{
		{
			UserName: "Anna Petrenko",
		},
	}

	free := [][]bool{
		{true, false, true},
		{false, false, true},
	}

	ts, _ := NewTeacherService(teachers)
	ts.SetBusyness(free)

	got := ts.GetAll()[0].Business

	// Перевірка: значення збігаються
	if !reflect.DeepEqual(got, free) {
		t.Errorf("SetBusyness failed: expected %v, got %v", free, got)
	}

	// Перевірка: глибока копія (зміна в `free` не впливає на ts.teachers)
	free[0][0] = false

	if got[0][0] == false {
		t.Errorf("SetBusyness should copy data, but got changed after modifying input")
	}
}
