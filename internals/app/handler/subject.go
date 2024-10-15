package handler

import (
	"errors"
	"fmt"
	"ggkit_learn_service/internals/app/processor"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type SubjectHandler struct {
	processor *processor.SubjectProcessor
}

func NewSubjectHandler(processor *processor.SubjectProcessor) *SubjectHandler {
	handler := new(SubjectHandler)
	handler.processor = processor
	return handler
}

func (handler *SubjectHandler) List(w http.ResponseWriter, r *http.Request) {

	list, err := handler.processor.SubjectsList()
	if err != nil {
		WrapError(w, err)
		return
	}
	var m = map[string]interface{}{
		"result": "ok",
		"data":   list,
	}
	WrapOK(w, m)

}

func (handler *SubjectHandler) One(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data, err := handler.processor.SubjectById(vars["id"])
	if err != nil {
		WrapError(w, err)
		return
	}
	var m = map[string]interface{}{
		"result": "OK",
		"data":   data,
	}
	WrapOK(w, m)

}

func (handler *SubjectHandler) UploadSubject(w http.ResponseWriter, r *http.Request) {
	// Получаем поля title и description из тела запроса
	title := r.FormValue("title")
	description := r.FormValue("description")

	// Получаем файл из поля image
	file, header, err := r.FormFile("image")
	if err != nil {
		WrapError(w, err)
		return
	}
	defer file.Close()

	// Проверяем размер файла
	if header.Size > 20*1024*1024 {
		WrapError(w, errors.New("pазмер изображения должен быть не более 20 МБ"))
		return
	}

	// Создаем файл в публичной папке /images
	image_url := "praxis" + header.Filename

	newSubjectId, err := handler.processor.UploadSubject(title, description, image_url)
	if err != nil {
		WrapError(w, err)
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
		"result":     fmt.Sprintf("Изображение успешно сохранено: %s\n", header.Filename),
		"subject_id": newSubjectId,
	}
	WrapOK(w, m)

}

func (handler *SubjectHandler) UpdateSubject(w http.ResponseWriter, r *http.Request) {
	// Получаем поля title и description из тела запроса
	subject_id := r.FormValue("subject_id")
	title := r.FormValue("title")
	description := r.FormValue("description")
	var image_url string
	var newSubjectId int
	// Получаем файл из поля image
	file, header, err := r.FormFile("image")
	if err != nil {
		// Создаем файл в публичной папке /images
		image_url = ""
		logrus.Error(err)
		// WrapError(w, err)
		// return
		newSubjectId, err = handler.processor.UpdateSubject(subject_id, title, description, image_url)
		if err != nil {
			WrapError(w, err)
		}
	} else {
		// Проверяем размер файла
		if header.Size > 20*1024*1024 {
			WrapError(w, errors.New("pазмер изображения должен быть не более 20 МБ"))
			return
		}
		// Создаем файл в публичной папке /images
		image_url = "praxis" + header.Filename
		defer file.Close()

		newSubjectId, err = handler.processor.UpdateSubject(subject_id, title, description, image_url)
		if err != nil {
			WrapError(w, err)
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

	}

	// Возвращаем успешный ответ

	var m = map[string]interface{}{
		"result":     fmt.Sprintf("Изображение успешно сохранено: %s\n", image_url),
		"subject_id": newSubjectId,
	}
	WrapOK(w, m)

}

func (handler *SubjectHandler) DeleteSubject(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	err := handler.processor.DeleteSubject(vars["id"])
	if err != nil {
		WrapError(w, err)
		return
	}
	var m = map[string]interface{}{
		"result": "OK",
	}
	WrapOK(w, m)

}

func (handler *SubjectHandler) Image(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	imageName := queryParams.Get("id")
	if imageName != "" {
		imagePath := "./images/" + imageName
		WrapOKImage(w, imagePath)
	}
	WrapError(w, errors.New("Имя изображения не указано"))
}
