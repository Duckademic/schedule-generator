package entities

type ScheduleResult struct {
	TeacherWindows                int // Sum of teachers windows
	TeacherLessonOverlapping      int // Sum of overlapping lessons
	StudentGroupWindows           int // Sum of student group windows
	StudentGroupLessonOverlapping int // Sum of overlapping lessons
	HoursDeficit                  int // Sum of pending hours
}

// The lower fault the better schedule (if 0 - this is ideal schedule).
// Time complexity O(1)
func (sr *ScheduleResult) Fault() (fault float64) {
	return float64(sr.TeacherLessonOverlapping)*100 +
		float64(sr.StudentGroupLessonOverlapping)*100 +
		float64(sr.TeacherWindows)*10 +
		float64(sr.StudentGroupWindows)*1000 +
		float64(sr.HoursDeficit)*100
}
