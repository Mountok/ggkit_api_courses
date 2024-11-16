package processor

import (
	"ggkit_learn_service/internals/app/db"
	"ggkit_learn_service/internals/app/models"
	"strconv"
)

type SubjectTestProcessor struct {
	storage *db.SubjectTestStorage
}

func NewSubjectTestProcessor(storage *db.SubjectTestStorage) *SubjectTestProcessor {
	processor := new(SubjectTestProcessor)
	processor.storage = storage
	return processor
}

// ---------------------- ТЕСТЫ

func (processor *SubjectTestProcessor) CreateTestForSubject(testTitle,subjectId string ) (int, error) {
	return processor.storage.CreateTestForSubject(testTitle,subjectId)
} 
func (processor *SubjectTestProcessor) ReadTestsForSubject(subjectId string) ([]models.SubjectTest,error) {
	return processor.storage.ReadTestsForSubject(subjectId)
}
func (processor *SubjectTestProcessor) UpdateTestForSubject(newTitle,subjectId,lastTitle string ) (int, error) {
	return processor.storage.UpdateTestForSubject(newTitle,subjectId,lastTitle)
} 
func (processor *SubjectTestProcessor) DeleteTestForSubject(testTitle,subjectId string ) error {
	return processor.storage.DeleteTestForSubject(testTitle,subjectId)
} 

// ------------------ ВОПРОСЫ К ТЕСТАМ

func (processor *SubjectTestProcessor) CreateQuestionForTest(testId,question,options,answer string) (int,error) {
	return processor.storage.CreateQuestionForTest(testId,question,options,answer)
}

func (processor *SubjectTestProcessor) ReadQuestionForTest(testId string) ([]models.TestQuestion, error) {
	return processor.storage.ReadQuestionForTest(testId)
}

func (processor *SubjectTestProcessor) UpdateQuestionForTest(testId,question,options,answer string) (int,error){
	return processor.storage.UpdateQuestionForTest(testId,question,options,answer)
}

func (processor *SubjectTestProcessor) DeleteQuestionForTest(testId,questionId string) error {
	return processor.storage.DeleteQuestionForTest(testId,questionId)
}
// ----------------- ВЫПОЛНЕНЫЕ ТЕСТЫ

func (processor *SubjectTestProcessor) CreateComletedTest(testId,userId, pointsString string) (int,error) {
	points, err := strconv.Atoi(pointsString)
	if err != nil {
		return 0,err
	}
	return processor.storage.CreateCompletedTest(userId,testId,points)
}

func (processor *SubjectTestProcessor) ReadComletedTest(subject_id,userId string) ([]models.CompletedTestCheck, error) {
	subjectIDInt, err := strconv.Atoi(subject_id)
   if err != nil {
       return []models.CompletedTestCheck{}, err
   }
	return processor.storage.ReadCompletedTest(subjectIDInt,userId)
}

func (processor *SubjectTestProcessor) UpdateComletedTest(testId,userId,pointsString string) (int,error){
	points, err := strconv.Atoi(pointsString)
	if err != nil {
		return 0,err
	}
	return processor.storage.UpdateCompletedTest(testId,userId,points)
}

func (processor *SubjectTestProcessor) DeleteCompletedTest(testId,userId,completedId string) error {
	return processor.storage.DeleteCompletedTest(testId,userId,completedId)
}





// --------------------- ПРОВЕРКА ВОПРОСОВ

func (processor *SubjectTestProcessor) CheckQuestion(answers []models.QuestionCheckReq,test_id,user_id,subject_id string) (int, error) {
	subjectIDInt, err := strconv.Atoi(subject_id)
	if err != nil {
		return 0, err
	}
	return processor.storage.CheckQuestion(answers,test_id,user_id,subjectIDInt)
}


// --------------------  ПОЛУЧЕНИЕ ВСЕХ ТЕСТОВ

func (processor *SubjectTestProcessor) GetAllCompleted(user_id string) (int, error) {
	return processor.storage.GetAllCompleted(user_id)
}