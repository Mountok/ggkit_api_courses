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
	return process.storage.CreateLesson(theme_id, theme_html)
}
func (process *LessonsProcessor) UpdateLesson(theme_id, theme_html string) error {
	if theme_id == "" || theme_id == " " {
		return errors.New("неверный theme_id")
	}
	if theme_html == " " || theme_html == "" {
		return errors.New("неверный theme_html")
	}
	return process.storage.UpdateLesson(theme_id, theme_html)
}
func (process *LessonsProcessor) GetLessonHTML(theme_id string) (string, error) {
	if theme_id == "" || theme_id == " " {
		return "", errors.New("неверный theme_id")
	}
	return process.storage.GetLessonHTML(theme_id)
}

func (process *LessonsProcessor) GetLessonByIdSubjectAndTheme(subjectId, themeId int) (models.LessonResponse, error) {
	if subjectId <= 0 {
		return models.LessonResponse{}, fmt.Errorf("uncorrect subject id = %d", subjectId)
	} else if themeId <= 0 {
		return models.LessonResponse{}, fmt.Errorf("uncorrect theme id = %d", themeId)
	}
	res := process.storage.GetLesson(subjectId, themeId)
	var lessons = models.LessonResponse{
		Id:      res[0].Id,
		Upkeep:  `` + res[0].Upkeep,
		ThemeId: res[0].ThemeId,
		Title:   res[0].Title,
	}
	fmt.Println(lessons)
	return lessons, nil
}
