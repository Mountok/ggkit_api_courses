package processor

import (
	"errors"
	"ggkit_learn_service/internals/app/db"
	"ggkit_learn_service/internals/app/models"
)

type UserProcessor struct {
	storage *db.UserStorage
}

func NewUserProcessor(storage *db.UserStorage) *UserProcessor {
	process := new(UserProcessor)
	process.storage = storage
	return process
}

func (processor *UserProcessor) LastSubject(userId string) ([]models.LastSubject, error) {
	return processor.storage.LastSubject(userId)
}
func (processor *UserProcessor) SetLastSubject(userId, subject_id string) error {
	if subject_id == "" || userId == "" {
		return errors.New("не верно переданы значения (subject_id/userid)")
	}
	return processor.storage.SetLastSubject(userId, subject_id)
}

func (processor *UserProcessor) UploadAvatar(userId, image_url string) error {
	return processor.storage.UploadAvatar(userId, image_url)
}

func (processor *UserProcessor) ChangeDescription(userId, newDescription string) error {
	return processor.storage.ChangeDescription(userId, newDescription)
}

func (processor *UserProcessor) ChangeUserName(userId, newName string) (string, error) {
	return processor.storage.ChangeUserName(userId, newName)
}

func (processor *UserProcessor) GetPoint(user_id, theme_id string) error {
	return processor.storage.GetPoint(user_id, theme_id)
}

func (processor *UserProcessor) Rating() ([]models.RatingUser, error) {
	return processor.storage.Rating()
}

func (processor *UserProcessor) CheckDoneLessons(user_id, theme_id string) (int, error) {
	return processor.storage.CheckDoneLessons(user_id, theme_id)
}

func (processor *UserProcessor) GetUserOnSubject(subjectId string) ([]models.FindUser, error) {
	if subjectId == "" || subjectId == "0" {
		return []models.FindUser{}, errors.New("undefined value for subject_id")
	}
	return processor.storage.GetUserOnSubject(subjectId)
}