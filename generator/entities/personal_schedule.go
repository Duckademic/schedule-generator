package entities

import (
	"fmt"
	"os"
)

type PersonalSchedule struct {
	BusyGrid *BusyGrid
	Lessons  []*Lesson
	Out      string // шлях до текстового файлу
}

func (ps *PersonalSchedule) InsertLesson(l *Lesson) {
	index := len(ps.Lessons)
	for i := range ps.Lessons {
		if ps.Lessons[i].After(l) {
			index = i
			break
		}
	}

	ps.Lessons = append(ps.Lessons[:index], append([]*Lesson{l}, ps.Lessons[index:]...)...)
}

func (ps *PersonalSchedule) WritePS(lessonToString func(*Lesson) string) error {
	file, err := os.Create(ps.Out)
	if err != nil {
		return err
	}
	defer file.Close()

	lessonIndex := 0
	for day := range ps.BusyGrid.Grid {
		dayStr := []string{"Неділя", "Понеділок", "Вівторок", "Середа", "Четвер", "П'ятниця", "Субота"}[day%7]
		_, err := file.WriteString(fmt.Sprintf("%s (день %d) \n", dayStr, day))
		if err != nil {
			return err
		}

		for slot := range ps.BusyGrid.Grid[day] {
			var lStr string
			currentSlot := LessonSlot{Day: day, Slot: slot}
			if len(ps.Lessons) != lessonIndex && ps.Lessons[lessonIndex].Slot == currentSlot {
				lStr = lessonToString(ps.Lessons[lessonIndex])
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
