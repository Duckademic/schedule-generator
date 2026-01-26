package entities

// LoadService tracks and evaluates the study workload.
type LoadService interface {
	AddLesson(lesson *Lesson)      // Registers a lesson.
	CountHourDeficit() int         // Returns the number of missing study hours.
	IsEnoughLessons() bool         // Returns true if the entity doesn't require additional lessons.
	GetAssignedLessons() []*Lesson // Returns registered lessons as an array.
}

// NewLoadService creates a new LoadChecker basic instance.
//
// It requires a count of required hours for the entity.
func NewLoadService(rh int) LoadService {
	return &loadService{requiredHours: rh}
}

// loadService is the basic implementation of the LoadChecker interface.
type loadService struct {
	lessons       []*Lesson
	requiredHours int
	currentHours  int
}

func (lc *loadService) AddLesson(lesson *Lesson) {
	lc.lessons = append(lc.lessons, lesson)
	lc.currentHours += lesson.Value
}
func (lc *loadService) CountHourDeficit() int {
	count := lc.requiredHours - lc.currentHours
	if count > 0 {
		return count
	}

	return 0
}
func (lc *loadService) IsEnoughLessons() bool {
	return lc.CountHourDeficit() <= 0
}
func (lc *loadService) GetAssignedLessons() []*Lesson {
	return lc.lessons
}
