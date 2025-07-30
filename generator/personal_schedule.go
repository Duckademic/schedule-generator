package generator

import (
	"fmt"
	"os"
)

type PersonalSchedule struct {
	busyGrid *BusyGrid
	lessons  []*Lesson
	out      string // шлях до текстового файлу
}

func (ps *PersonalSchedule) InsertLesson(l *Lesson) {
	index := len(ps.lessons)
	for i := range ps.lessons {
		if ps.lessons[i].After(l) {
			index = i
			break
		}
	}

	ps.lessons = append(ps.lessons[:index], append([]*Lesson{l}, ps.lessons[index:]...)...)
}

func (ps *PersonalSchedule) WritePS(lessonToString func(*Lesson) string) error {
	file, err := os.Create(ps.out)
	if err != nil {
		return err
	}
	defer file.Close()

	lessonIndex := 0
	for day := range ps.busyGrid.Grid {
		dayStr := []string{"Неділя", "Понеділок", "Вівторок", "Середа", "Четвер", "П'ятниця", "Субота"}[day%7]
		_, err := file.WriteString(fmt.Sprintf("%s (день %d) \n", dayStr, day))
		if err != nil {
			return err
		}

		for slot := range ps.busyGrid.Grid[day] {
			var lStr string
			currentSlot := LessonSlot{Day: day, Slot: slot}
			if len(ps.lessons) != lessonIndex && ps.lessons[lessonIndex].Slot == currentSlot {
				lStr = lessonToString(ps.lessons[lessonIndex])
				lessonIndex++
			}

			_, err := file.WriteString(fmt.Sprintf("%d. %s\n", slot+1, lStr))
			if err != nil {
				return err
			}
		}

		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}
