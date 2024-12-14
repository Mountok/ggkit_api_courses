package handler

import (
	"errors"
	"fmt"
	"ggkit_learn_service/internals/app/processor"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	processor *processor.UserProcessor
}

func NewUserHandler(processor *processor.UserProcessor) *UserHandler {
	handler := new(UserHandler)
	handler.processor = processor
	return handler
}

func (handler *UserHandler) LastSubject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
	userIdString := w.Header().Get(UserCtx)
	switch r.Method {
	case http.MethodGet:
		logrus.Print(r.Method)
		subjectArray, err := handler.processor.LastSubject(userIdString)
		if err != nil {
			WrapError(w, err)
			return
		}
		var m = map[string]interface{}{
			"result": "OK",
			"data":   subjectArray,
		}
		WrapOK(w, m)

	case http.MethodPost:
		courseId := vars["course_id"]

		err := handler.processor.SetLastSubject(userIdString, courseId)
		if err != nil {
			WrapError(w, err)
			return
		}
		var m = map[string]interface{}{
			"result": "OK",
			"data":   courseId,
		}
		WrapOK(w, m)
	}
}

func (handler *UserHandler) ChangeName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var (
		userIdString = vars["user_id"]
		newName      = vars["new_name"]
	)

	newName, err := handler.processor.ChangeUserName(userIdString, newName)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result": http.StatusOK,
		"data":   newName,
	}
	WrapOK(w, m)
}

func (handler *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	// Получаем поля title и description из тела запроса
	userId := r.FormValue("user_id")

	// Получаем файл из поля image
	file, header, err := r.FormFile("image")
	if err != nil {
		WrapError(w, err)
		return
	}
	defer file.Close()

	// Проверяем размер файла
	if header.Size > 10*1024*1024 {
		WrapError(w, errors.New("pазмер изображения должен быть не более 10 МБ"))
		return
	}

	// Создаем файл в публичной папке /images
	image_url := "avatar_for_user_" + userId + ".webp"

	err = handler.processor.UploadAvatar(userId, image_url)
	if err != nil {
		WrapError(w, err)
		return
	}

	out, err := os.Create("./images/" + image_url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Копируем содержимое файла в созданный файл
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ

	var m = map[string]interface{}{
		"result":  fmt.Sprintf("Изображение успешно сохранено: %s\n", header.Filename),
		"user_id": userId,
	}
	WrapOK(w, m)

}

func (handler *UserHandler) ChangeDescription(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("user_id")
	newDescription := r.FormValue("description")

	err := handler.processor.ChangeDescription(userId, newDescription)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]interface{}{
		"result":  "Описание провеля успешно заменено",
		"user_id": userId,
	}

	WrapOK(w, m)

}

func (handler *UserHandler) GetPoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w,r,err := UserIdentify(w,r)
	if err != nil {
		WrapErrorWithStatus(w,err,http.StatusUnauthorized)
		return
	}

	user_id := w.Header().Get(UserCtx)
	theme_id := vars["theme_id"]

	record, err := handler.processor.CheckDoneLessons(user_id, theme_id)
	if err != nil {
		WrapError(w, err)
		return
	}
	if record != 0 {
		var m = map[string]interface{}{
			"result":    "NOT UPDATE",
			"user_id":   user_id,
			"lesson_id": theme_id,
		}
		WrapOK(w, m)
	} else {
		err := handler.processor.GetPoint(user_id, theme_id)
		if err != nil {
			WrapError(w, err)
			return
		}
		var m = map[string]interface{}{
			"result":    "OK",
			"user_id":   user_id,
			"lesson_id": theme_id,
		}
		WrapOK(w, m)
	}
}

func (handler *UserHandler) Rating(w http.ResponseWriter, r *http.Request) {
	lists, err := handler.processor.Rating()
	if err != nil {
		WrapError(w, err)
		return
	}
	var m = map[string]interface{}{
		"result": "OK",
		"data":   lists,
	}
	WrapOK(w, m)
}
