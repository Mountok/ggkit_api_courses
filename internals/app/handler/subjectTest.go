package handler

import (
	"encoding/json"
	"fmt"
	"ggkit_learn_service/internals/app/models"
	"ggkit_learn_service/internals/app/processor"
	"net/http"

	"github.com/gorilla/mux"
)

type SubjectTestHandler struct {
	processor *processor.SubjectTestProcessor
}

func NewSubjectTestHandler(processor *processor.SubjectTestProcessor) *SubjectTestHandler {
	handler := new(SubjectTestHandler)
	handler.processor = processor
	return handler
}

func (handler *SubjectTestHandler) TestsForSubject(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)
		m = map[string]interface{}{}
		subjectId = vars["subject_id"]
	)
	switch r.Method {
	case http.MethodPost:
		testTitle := r.FormValue("title")
		id, err := handler.processor.CreateTestForSubject(testTitle, subjectId)
		if err != nil {
			WrapError(w, err)
			return
		}
		m = map[string]interface{}{
			"result":  "OK",
			"title":   testTitle,
			"test_id": id,
		}
	case http.MethodGet:
		listsOfTests, err := handler.processor.ReadTestsForSubject(subjectId)
		if err != nil {
			WrapError(w, err)
			return
		}
		m = map[string]interface{}{
			"result": "OK",
			"data":   listsOfTests,
		}
	case http.MethodPut:
		newTitle := r.FormValue("title")
		lastTitle := r.FormValue("last_title")
		id, err := handler.processor.UpdateTestForSubject(newTitle, subjectId, lastTitle)
		if err != nil {
			WrapError(w, err)
			return
		}
		m = map[string]interface{}{
			"result":  "OK",
			"test_id": id,
			"data":    fmt.Sprintf("Тест '%s' переименован на '%s' для предмета с id %s", lastTitle, newTitle, subjectId),
		}
	case http.MethodDelete:
		testTitle := r.FormValue("title")
		err := handler.processor.DeleteTestForSubject(testTitle, subjectId)
		if err != nil {
			WrapError(w, err)
		}
		m = map[string]interface{}{
			"result": "ok",
			"data":   fmt.Sprintf("Тест '%s' для курса с id %s - успешно удален.", testTitle, subjectId),
		}
	}
	WrapOK(w, m)

}

func (handler *SubjectTestHandler) TestsQuestions(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)
		m      = map[string]interface{}{}
		testId = vars["test_id"]
	)
	switch r.Method {
	case http.MethodPost:
		var (
			question = r.FormValue("question")
			options  = r.FormValue("options")
			answer   = r.FormValue("answer")
		)
		questionId, err := handler.processor.CreateQuestionForTest(testId, question, options, answer)
		if err != nil {
			WrapError(w,err)
			return
		}
		m = map[string]interface{}{
			"result": "OK",
			"data": questionId,
		}
	case http.MethodGet:
		questionList, err := handler.processor.ReadQuestionForTest(testId);
		if err != nil {
			WrapError(w,err)
			return
		}
		m = map[string]interface{}{
			"result": "OK",
			"data": questionList,
		}
	case http.MethodPut:
		var (
			question = r.FormValue("question")
			options  = r.FormValue("options")
			answer   = r.FormValue("answer")
		)
		updateQuestionId, err := handler.processor.UpdateQuestionForTest(testId,question,options,answer)
		if err != nil {
			WrapError(w,err)
			return
		}
		m = map[string]interface{}{
			"result": "OK",
			"updated_question_id": updateQuestionId,
		}
	case http.MethodDelete:
		questionId := r.FormValue("question_id")
		err := handler.processor.DeleteQuestionForTest(testId,questionId)
		if err != nil {
			WrapError(w,err)
			return
		}
		m = map[string]interface{}{
			"result": "OK",
			"data": "Вопрос для теста удален",
		}
	}

	
	WrapOK(w, m)
}

func (handler *SubjectTestHandler) CompletedTest(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)
		testId = vars["test_id"]
		userId = r.FormValue("user_id")
		m = make(map[string]interface{})
	)
	switch r.Method {
	case http.MethodPost:
		pointsString := r.FormValue("points")
		compledetTestId, err := handler.processor.CreateComletedTest(testId,userId,pointsString)
		if err != nil {
			WrapError(w,err)
			return
		}
		m = map[string]interface{}{
			"result": "OK",
			"completed_test_id": compledetTestId,
		}
	case http.MethodGet:
		completedTest, err := handler.processor.ReadComletedTest(testId,userId)
		if err != nil {
			WrapError(w,err)
			return
		}
		m = map[string]interface{}{
			"result": "OK",
			"data": completedTest,
		}
	case http.MethodPut:
		pointsString := r.FormValue("points")
		testId, err := handler.processor.UpdateComletedTest(testId,userId,pointsString)
		if err != nil {
			WrapError(w,err)
			return
		}
		m = map[string]interface{}{
			"result": "OK",
			"updated_completed_test_id": testId,
		}
		
	case http.MethodDelete:
		completedId := r.FormValue("completed_test_id")
		err := handler.processor.DeleteCompletedTest(testId,userId,completedId)
		if err != nil {
			WrapError(w,err)
			return
		}
		m = map[string]interface{}{
			"result": "OK",
			"data": fmt.Sprintf("Удален пройденный курс с id=%s",testId),
		}
	}

	WrapOK(w,m)

}

func (handler *SubjectTestHandler) CheckQuestion(w http.ResponseWriter, r *http.Request) {
	var resp []models.QuestionCheckReq
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		WrapError(w,err)
	}
	points, err := handler.processor.CheckQuestion(resp)
	if err != nil {
		WrapError(w,err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"points": points,
	}
	WrapOK(w,m)
}