package handler

import (
	"fmt"
	"ggkit_learn_service/internals/app/processor"
	"net/http"

	"github.com/gorilla/mux"
)

type ThemesHandler struct {
	processor *processor.ThemesProcessor
}

func NewThemesHandler(processor *processor.ThemesProcessor) *ThemesHandler {
	handler := new(ThemesHandler)
	handler.processor = processor
	return handler
}

func (handler *ThemesHandler) GetAllCompleted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id := vars["user_id"]
	subject_id := vars["subject_id"]

	if user_id == "" || subject_id == "" {
		WrapError(w, fmt.Errorf("not valid (user_id or subject_id)"))
		return
	}

	listsOfCompletedThemes, err := handler.processor.GetAllCompeted(user_id, subject_id)

	if err != nil {
		WrapError(w, err)
		return
	}
	var m = map[string]interface{}{
		"results": "ok",
		"data":    listsOfCompletedThemes,
	}
	WrapOK(w, m)
}

func (handler *ThemesHandler) CreateTheme(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	subject_id := r.FormValue("subject_id")

	newThemeId, err := handler.processor.CreateTheme(title, description, subject_id)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result":     "OK",
		"theme_id":   newThemeId,
		"subject_id": subject_id,
	}

	WrapOK(w, m)

}

func (handler *ThemesHandler) Themes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data, err := handler.processor.ThemesBySubjectId(vars)
	if err != nil {
		WrapError(w, err)
	}
	var m = map[string]interface{}{
		"result": "OK",
		"data":   data,
	}
	WrapOK(w, m)
}
