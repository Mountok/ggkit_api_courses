package handler

import (
	"fmt"
	"ggkit_learn_service/internals/app/processor"
	"net/http"

	// "strconv"

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

func (handler *ThemesHandler) GetAllCompletedBySubject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
	user_id := w.Header().Get(UserCtx)
	subject_id := vars["subject_id"]

	if user_id == "" || subject_id == "" {
		WrapError(w, fmt.Errorf("not valid (user_id or subject_id)"))
		return
	}

	listsOfCompletedThemes, err := handler.processor.GetAllCompletedBySubject(user_id, subject_id)

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

func (handler *ThemesHandler) GetAllCompleted(w http.ResponseWriter, r *http.Request) {
	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
	user_id := w.Header().Get(UserCtx)

	if user_id == "" {
		WrapError(w, fmt.Errorf("not valid user_id"))
		return
	}

	listsOfCompletedThemes, err := handler.processor.GetAllCompeted(user_id)

	if err != nil {
		WrapError(w, err)
		return
	}
	var m = map[string]interface{}{
		"results": "ok",
		"data":    len(listsOfCompletedThemes),
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

func (handler *ThemesHandler) UpdateTheme(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	themeIdStrin := vars["theme_id"]
	themeTitle := r.FormValue("theme_title")
	themeDescription := r.FormValue("theme_description")
	id, err := handler.processor.UpdateTheme(themeIdStrin, themeTitle, themeDescription)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": "OK",
		"data":   id,
	}

	WrapOK(w, m)

}

func (handler *ThemesHandler) DeleteTheme(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	theme_id := vars["theme_id"]
	err := handler.processor.DeleteTheme(theme_id)
	if err != nil {
		WrapError(w, err)
		return
	}
	var m = map[string]interface{}{
		"result": "OK",
		"data":   fmt.Sprintf("тема с id=%s удалена", theme_id),
	}

	WrapOK(w, m)
}

func (handler *ThemesHandler) Themes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
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
