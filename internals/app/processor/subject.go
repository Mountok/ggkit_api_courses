package processor

import (
	"errors"
	"fmt"
	"ggkit_learn_service/internals/app/db"
	"ggkit_learn_service/internals/app/models"
	"strconv"
)

type SubjectProcessor struct {
	storage *db.SubjectStorage
}

func NewSubjectProcessor(cache *db.SubjectStorage) *SubjectProcessor {
	processor := new(SubjectProcessor)
	processor.storage = cache
	return processor
}

func (process *SubjectProcessor) SubjectsList(userId string) ([]models.SubjectResponse, error) {
	return process.storage.GetAllSubjects(userId)
}
func (process *SubjectProcessor) SubjectById(id string) ([]models.Subject, error) {
	num, err := strconv.Atoi(id)
	if err != nil {
		return []models.Subject{}, fmt.Errorf("uncorrect id - (%s) is not integer", id)
	}
	if num <= 0 {
		return []models.Subject{}, fmt.Errorf("uncorrect id (%s)", id)
	}
	return process.storage.GetSubjectById(num)
}

func (process *SubjectProcessor) UploadSubject(title, description, image_url, is_certificated string) (int, error) {
	if title == "" || title == " " {
		return 0, errors.New("ошибка: заголовок предмета задан не верно")
	}
	return process.storage.UploadStorage(title, description, image_url, is_certificated)
}

func (process *SubjectProcessor) UpdateSubject(subject_id, title, description, image_url, is_certificated string) (int, error) {
	return process.storage.UpdateSubject(subject_id, title, description, image_url, is_certificated)
}
func (process *SubjectProcessor) DeleteSubject(id string) error {
	return process.storage.DeleteSubject(id)
}
func (process *SubjectProcessor) GetDeletedSubject() ([]int, error) {
	return process.storage.GetDeletedSubject()
}

func (process *SubjectProcessor) Certificate(subjectId, userId string) (interface{}, error) {
	return process.storage.Certificate(subjectId, userId)
}
