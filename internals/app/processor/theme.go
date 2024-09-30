package processor

import (
	"errors"
	"fmt"
	"ggkit_learn_service/internals/app/db"
	"ggkit_learn_service/internals/app/models"
	"strconv"
)

type ThemesProcessor struct {
	storage *db.ThemesStorage
}

func NewThemesProcessor(storage *db.ThemesStorage) *ThemesProcessor {
	processor := new(ThemesProcessor)
	processor.storage = storage
	return processor
}

func (process *ThemesProcessor) CreateTheme(title, description, subject_id string) (int, error) {
	if title == "" || title == " " {
		return 0, errors.New("неверно передано значение title")
	}
	if subject_id == "" || subject_id == " " {
		return 0, errors.New("неверно передано значение subject_id")
	}
	return process.storage.CreateTheme(title, description, subject_id)

}

func (process *ThemesProcessor) ThemesBySubjectId(req_vars map[string]string) ([]models.Theme, error) {
	num, err := strconv.Atoi(req_vars["subject_id"])
	if err != nil {
		return []models.Theme{}, fmt.Errorf("error: %s", err.Error())
	}
	return process.storage.GetThemesBySubjectId(num), nil

}

func (process *ThemesProcessor) GetAllCompeted(user_id, subject_id string) ([]int, error) {
	return process.storage.GetAllCompleted(user_id, subject_id)
}
