package processor

import (
	"errors"
	"fmt"
	"ggkit_learn_service/internals/app/db"
	"ggkit_learn_service/internals/app/models"
)

type LessonsProcessor struct {
	storage *db.LessonsStorage
}

func NewLessonsProcessor(storage *db.LessonsStorage) *LessonsProcessor {
	processor := new(LessonsProcessor)
	processor.storage = storage
	return processor
}

func (process *LessonsProcessor) CreateLesson(theme_id, theme_html string) error {
	if theme_id == "" || theme_id == " " {
		return errors.New("неверный theme_id")
	}
	if theme_html == " " || theme_html == "" {
		return errors.New("неверный theme_html")
	}
	return process.storage.CreateLesson(theme_id,theme_html)
}

func (process *LessonsProcessor) GetLessonByIdSubjectAndTheme(subjectId, themeId int) (models.Lesson, error) {
	if subjectId <= 0 {
		return models.Lesson{}, fmt.Errorf("uncorrect subject id = %d", subjectId)
	} else if themeId <= 0 {
		return models.Lesson{}, fmt.Errorf("uncorrect theme id = %d", themeId)
	}
	res := process.storage.GetLesson(subjectId, themeId)
	var lessons = models.Lesson{
		Id:      res[0].Id,
		Upkeep:  `` + res[0].Upkeep,
		ThemeId: res[0].ThemeId,
	}
	fmt.Println(lessons)
	return lessons, nil
}
