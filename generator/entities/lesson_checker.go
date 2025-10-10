package entities

type LessonChecker struct {
	Lessons       []*Lesson
	RequiredHours int
	CurrentHours  int
}

func (lc *LessonChecker) AddLesson(lesson *Lesson) {
	lc.Lessons = append(lc.Lessons, lesson)
	lc.CurrentHours += lesson.Value
}

// func (lc *LessonChecker) FullCheck() error {
// 	err := lc.CheckHours()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (lc *LessonChecker) CheckHours() error {
// 	if lc.RequiredHours > lc.CurrentHours {
// 		return fmt.Errorf("not enough hours (get %d, need %d)", lc.CurrentHours, lc.RequiredHours)
// 	}

// 	return nil
// }

// Counts hour deficit (0 is minimum).
// Time complexity O(1)
func (lc *LessonChecker) CountHourDeficit() int {
	count := lc.RequiredHours - lc.CurrentHours
	if count > 0 {
		return count
	}

	return 0
}
