package handler

import (
	"fmt"
	"ggkit_learn_service/internals/app/processor"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type LessonsHandler struct {
	processor *processor.LessonsProcessor
}

func NewLessonsHanler(processor *processor.LessonsProcessor) *LessonsHandler {
	handler := new(LessonsHandler)
	handler.processor = processor
	return handler
}

func (handler *LessonsHandler) CreateLesson(w http.ResponseWriter, r *http.Request) {
	var (
		theme_id   = r.FormValue("theme_id")
		theme_html = r.FormValue("theme_html")
	)

	err := handler.processor.CreateLesson(theme_id, theme_html)
	if err != nil {
		WrapError(w, err)
		return
	}
	var m = map[string]interface{}{
		"result": "ok",
	}
	WrapOK(w, m)

}

func (handler *LessonsHandler) Lesson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectId, err := strconv.Atoi(vars["subject_id"])
	if err != nil {
		WrapError(w, fmt.Errorf("subject id is not a number: %v", err))
	}
	themeId, err := strconv.Atoi(vars["theme_id"])
	if err != nil {
		WrapError(w, fmt.Errorf("subject id is not a number: %v", err))
	}
	data, err := handler.processor.GetLessonByIdSubjectAndTheme(subjectId, themeId)
	if err != nil {
		WrapError(w, err)
	}
	var m = map[string]interface{}{
		"result": "OK",
		"data":   data,
	}
	WrapOK(w, m)
}
