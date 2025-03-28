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

	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
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
	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
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
	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
	// Получаем поля title и description из тела запроса
	title := r.FormValue("title")
	description := r.FormValue("description")
	is_certificated := r.FormValue("iscertificated")

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

	image_url := "praxis_course_of_id_"

	newSubjectId, err := handler.processor.UploadSubject(title, description, image_url, is_certificated)
	if err != nil {
		WrapError(w, err)
	}
	// Создаем файл в публичной папке /images
	image_url = fmt.Sprintf("praxis_course_of_id_%d.webp", newSubjectId)
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
	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
	// Получаем поля title и description из тела запроса
	subject_id := r.FormValue("subject_id")
	title := r.FormValue("title")
	description := r.FormValue("description")
	is_certificated := r.FormValue("iscertificated")

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
		newSubjectId, err = handler.processor.UpdateSubject(subject_id, title, description, image_url, is_certificated)
		if err != nil {
			WrapError(w, err)
			return
		}
	} else {
		// Проверяем размер файла
		if header.Size > 20*1024*1024 {
			WrapError(w, errors.New("pазмер изображения должен быть не более 20 МБ"))
			return
		}
		// Создаем файл в публичной папке /images
		image_url = "praxis_course_of_id_" + subject_id + ".webp"
		defer file.Close()

		newSubjectId, err = handler.processor.UpdateSubject(subject_id, title, description, image_url, is_certificated)
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
	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
	subjectID := r.URL.Query().Get("id")
	switch r.Method{
	case http.MethodDelete:
		err = handler.processor.DeleteSubject(subjectID)
		if err != nil {
			WrapError(w, err)
			return
		}
		var m = map[string]interface{}{
			"result": "OK",
		}
		WrapOK(w, m)
	case http.MethodGet:
		ids, err := handler.processor.GetDeletedSubject()
		if err != nil {
			WrapError(w, err)
			return
		}
		var m = map[string]interface{}{
			"result": "OK",
			"data": ids,
		}
		WrapOK(w, m)
	}
	

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

func (handler *SubjectHandler) Video(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	videoName := queryParams.Get("id")
	if videoName == "" {
		WrapError(w, errors.New("Имя видео не указано"))
		return
	}

	videoPath := "./videos/" + videoName

	// Проверяем, существует ли файл
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		WrapError(w, fmt.Errorf("файл не найден"))
		return
	}

	// Отправляем файл с поддержкой потокового воспроизведения
	w.Header().Set("Content-Type", "video/mp4") // Укажи нужный формат
	w.Header().Set("Accept-Ranges", "bytes")    // Включаем поддержку Range-запросов
	http.ServeFile(w, r, videoPath)
}
func (handler *SubjectHandler) Certificate(w http.ResponseWriter, r *http.Request) {
	w, r, err := UserIdentify(w, r)
	if err != nil {
		WrapErrorWithStatus(w, err, http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	subjectId := vars["subject_id"]
	userId := w.Header().Get(UserCtx)
	var m = map[string]interface{}{}
	switch r.Method {
	case http.MethodGet:
		date, err := handler.processor.Certificate(subjectId, userId)
		if err != nil {
			WrapError(w, err)
			return
		}
		m = map[string]interface{}{
			"result":     "OK",
			"courseDone": true,
			"date": date,
		}
	}

	WrapOK(w, m)

}

