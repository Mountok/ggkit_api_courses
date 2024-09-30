package processor

import (
	"ggkit_learn_service/internals/app/models"
	"ggkit_learn_service/internals/app/rdb"
)

type SubjectProcessor struct {
	cache *rdb.SubjectCache
}

func NewSubjectProcessor(cache *rdb.SubjectCache) *SubjectProcessor {
	processor := new(SubjectProcessor)
	processor.cache = cache
	return processor
}

func (process *SubjectProcessor) SubjectsList() ([]models.Subject, error) {
	return process.cache.GetAllSubjects()
}
func (process *SubjectProcessor) SubjectById(id string) ([]models.Subject, error) {
	// num, err := strconv.Atoi(id)
	// if err != nil {
	// 	return []models.Subject{}, fmt.Errorf("uncorrect id - (%s) is not integer", id)
	// }
	// if num <= 0 {
	// 	return []models.Subject{}, fmt.Errorf("uncorrect id (%s)", id)
	// }
	// return process.storage.GetSubjectById(num), nil
	return []models.Subject{},nil
}

func (process *SubjectProcessor) UploadSubject(title, description,image_url string) (int, error) {
	// if title == "" || title == " " {
	// 	return 0, errors.New("ошибка: заголовок предмета задан не верно")
	// }
	// return process.storage.UploadStorage(title, description,image_url)
	return 0,nil
}
