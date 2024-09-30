package api

import (
	"ggkit_learn_service/internals/app/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateRoute(
	subjectHandler *handler.SubjectHandler,
	themeHandler *handler.ThemesHandler,
	lessonsHandler *handler.LessonsHandler,
	loginHandler *handler.LoginHandler,
	userHandler *handler.UserHandler,
	testHandler *handler.SubjectTestHandler,
) *mux.Router {
	router := mux.NewRouter()

	// ! Эндпоинты для тем
	
	router.HandleFunc("/api/subject", subjectHandler.List).Methods(http.MethodGet)
	router.HandleFunc("/api/subject/{id}", subjectHandler.One).Methods(http.MethodGet)
	router.HandleFunc("/api/subject", subjectHandler.UploadSubject).Methods(http.MethodPost)
	////////////////////////////////////////////////////////
	// ! Эндпоинты для тем

	router.HandleFunc("/api/themes/{subject_id}", themeHandler.Themes).Methods(http.MethodGet)
	router.HandleFunc("/api/themes", themeHandler.CreateTheme).Methods(http.MethodPost)

	// TODO: Сделать эндпоинт для отображения всех пройденных уроков по id предмета
	router.HandleFunc("/api/themes/complete/{user_id}/{subject_id}",
		themeHandler.GetAllCompleted).Methods(http.MethodGet)

	////////////////////////////////////////////////////////
	// ! Эндпоинты для уроков

	router.HandleFunc("/api/lessons/{subject_id}/{theme_id}", lessonsHandler.Lesson).Methods(http.MethodGet)
	router.HandleFunc("/api/lessons", lessonsHandler.CreateLesson).Methods(http.MethodPost)

	////////////////////////////////////////////////////////
	// ! Авторизация / регистрация
	router.HandleFunc("/api/reg", loginHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/api/auth", loginHandler.Auth).Methods(http.MethodPost)
	router.HandleFunc("/api/validate", loginHandler.Validate).Methods(http.MethodPost)

	////////////////////////////////////////////////////////
	// ! Эндпоинты для профиля

	// * Получение данных (профиля) пользователя
	router.HandleFunc("/api/profile/{user_id}", loginHandler.Profile).Methods(http.MethodGet)
	// * Смена автарки
	router.HandleFunc("/api/profile/avatar", userHandler.UploadAvatar).Methods(http.MethodPost)
	// * Изменение описания
	router.HandleFunc("/api/profile/description", userHandler.ChangeDescription).Methods(http.MethodPost)
	// * Добавление поинтов (100 поинтов == 1 лвл)
	router.HandleFunc("/api/profile/point/{user_id}/{theme_id}", userHandler.GetPoint).Methods(http.MethodPost)
	// * Смена имени
	router.HandleFunc("/api/profile/name/{user_id}/{new_name}", userHandler.ChangeName).Methods(http.MethodPost)
	// * Получение рейтинга пользователя
	router.HandleFunc("/api/profiles", userHandler.Rating).Methods(http.MethodGet)

	// * Для отметки пройденной темы для пользователя по id
	router.HandleFunc("/api/profile/subject/{user_id}", userHandler.LastSubject).Methods(http.MethodGet)
	// * Получения всех пройденных тем по id пользователя и id курса
	router.HandleFunc("/api/profile/subject/{user_id}/{course_id}", userHandler.LastSubject).Methods(http.MethodPost)

	////////////////////////////////////////////////////////
	// ! Эндпоинты для тестов
	router.HandleFunc("/api/test/{subject_id}", testHandler.TestsForSubject).Methods(
		http.MethodGet,
		http.MethodPost,
		http.MethodDelete,
		http.MethodPut,
	)

	router.HandleFunc("/api/testing/{test_id}", testHandler.TestsQuestions).Methods(
		http.MethodGet,
		http.MethodPost,
		http.MethodDelete,
		http.MethodPut,
	)

	router.HandleFunc("/api/test/check/{question_id}", testHandler.CheckQuestion).Methods(
		http.MethodPost,
	)

	router.HandleFunc("/api/test/result/{test_id}", testHandler.CompletedTest).Methods(
		http.MethodGet,
		http.MethodPost,
		http.MethodDelete,
		http.MethodPut,
	)

	// ! Эндпоинт для публичного доступа к изображениям)
	// ! НЕ ТРОГАТЬ Я НЕ ЗНАЮ КАК НО ОНО РАБОТАЕТ
	router.HandleFunc("/images", subjectHandler.Image).Methods(http.MethodGet)

	return router
}
