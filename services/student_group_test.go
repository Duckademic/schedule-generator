package services

import (
	"testing"

	"github.com/Duckademic/schedule-generator/types"
)

func getTestGroups() []types.StudentGroup {
	return []types.StudentGroup{
		{Name: "GroupA"},
		{Name: "GroupB"},
		{Name: "GroupC"},
		{Name: "GroupD"},
		{Name: "GroupE"},
		{Name: "GroupF"},
	}
}

func TestFindGroup(t *testing.T) {
	groups := getTestGroups()
	service, _ := NewStudentGroupService(groups, 4)

	group := service.Find("GroupC")
	if group == nil || group.Name != "GroupC" {
		t.Errorf("expected to find GroupC, got %v", group)
	}

	group = service.Find("Unknown")
	if group != nil {
		t.Errorf("expected nil for unknown group, got %v", group)
	}
}

func TestCountGroupLessonsOn(t *testing.T) {
	groups := getTestGroups()
	groups[0].Business = [][]bool{
		{true, false, true},   // день 0: 2 заняття
		{false, false, false}, // день 1: 0 занять
	}
	service, _ := NewStudentGroupService(groups, 4)

	count := service.CountLessonsOn("GroupA", 0)
	if count != 2 {
		t.Errorf("expected 2 lessons on day 0, got %d", count)
	}

	count = service.CountLessonsOn("GroupA", 1)
	if count != 0 {
		t.Errorf("expected 0 lessons on day 1, got %d", count)
	}
}

func TestGetGroupLectureDay(t *testing.T) {
	groups := getTestGroups()

	// GroupA має пари тільки на дні 0 та 1 (2 заняття вже є в день 1)
	groups[0].Business = [][]bool{
		{false, false}, // день 0
		{true, false},  // день 1
		{false, false}, // день 2
	}

	groups[0].LectureDays = []int{1, 2} // допустимі дні

	service, _ := NewStudentGroupService(groups, 1)

	day := service.GetLectureDay("GroupA", 0)
	if day != 2 {
		t.Errorf("expected day 2 (free and allowed), got %d", day)
	}
}

func TestSetBusynessGroup(t *testing.T) {
	groups := getTestGroups()
	free := [][]bool{
		{true, false},
		{false, true},
	}

	service, _ := NewStudentGroupService(groups, 3)
	service.SetBusyness(free)

	for _, g := range service.GetAll() {
		if len(g.Business) != len(free) || len(g.Business[0]) != len(free[0]) {
			t.Errorf("incorrect business size after SetBusyness for group %s", g.Name)
		}
		if g.Business[0][0] != true || g.Business[1][1] != true {
			t.Errorf("expected copied business matrix, got incorrect values in group %s", g.Name)
		}
	}
}
