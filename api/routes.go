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
	router.HandleFunc("/api/subject", subjectHandler.UpdateSubject).Methods(http.MethodPut)
	router.HandleFunc("/api/delete/subject", subjectHandler.DeleteSubject).Methods(http.MethodDelete,http.MethodGet)
	////////////////////////////////////////////////////////
	// ! Эндпоинты для тем

	router.HandleFunc("/api/themes/{subject_id}", themeHandler.Themes).Methods(http.MethodGet)
	router.HandleFunc("/api/themes", themeHandler.CreateTheme).Methods(http.MethodPost)
	router.HandleFunc("/api/themes/{theme_id}", themeHandler.UpdateTheme).Methods(http.MethodPost)
	router.HandleFunc("/api/themes/{theme_id}", themeHandler.DeleteTheme).Methods(http.MethodDelete)

	// ! Поиск пройденных тем для пользователя по предмету
	router.HandleFunc("/api/themes/complete/{subject_id}",
		themeHandler.GetAllCompletedBySubject).Methods(http.MethodGet)
	
	// ! Поиск пройденных тем для пользователя общее количество
	router.HandleFunc("/api/themes/completed/{user_id}",
		themeHandler.GetAllCompleted).Methods(http.MethodGet)

	////////////////////////////////////////////////////////
	// ! Эндпоинты для уроков

	router.HandleFunc("/api/lessons/{subject_id}/{theme_id}", lessonsHandler.Lesson).Methods(http.MethodGet)

	router.HandleFunc("/api/lessons", lessonsHandler.GRUDLesson).Methods(
		http.MethodPost,
		http.MethodGet,
		http.MethodPut,
	)

	////////////////////////////////////////////////////////
	// ! Авторизация / регистрация
	
	router.HandleFunc("/api/sign-up", loginHandler.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/api/sign-in", loginHandler.SignIn).Methods(http.MethodPost)
	router.HandleFunc("/api/authorization",loginHandler.Authorization).Methods(http.MethodPost)
	// СТАРОЕ !!!!!
	// router.HandleFunc("/api/reg", loginHandler.Create).Methods(http.MethodPost)
	// router.HandleFunc("/api/auth", loginHandler.Auth).Methods(http.MethodPost)
	// router.HandleFunc("/api/validate", loginHandler.Validate).Methods(http.MethodPost)

	////////////////////////////////////////////////////////
	// ! Эндпоинты для профиля

	// * Получение данных (профиля) пользователя
	router.HandleFunc("/api/profile", loginHandler.Profile).Methods(http.MethodGet)
	// * Смена автарки
	router.HandleFunc("/api/profile/avatar", userHandler.UploadAvatar).Methods(http.MethodPost)
	// * Изменение описания
	router.HandleFunc("/api/profile/description", userHandler.ChangeDescription).Methods(http.MethodPost)
	// * Добавление поинтов (100 поинтов == 1 лвл)
	router.HandleFunc("/api/profile/point/{theme_id}", userHandler.GetPoint).Methods(http.MethodPost)
	// * Смена имени
	router.HandleFunc("/api/profile/name/{new_name}", userHandler.ChangeName).Methods(http.MethodPost)
	// * Получение рейтинга пользователя
	router.HandleFunc("/api/profiles", userHandler.Rating).Methods(http.MethodGet)

	// *  пройденного предмета для пользователя по id
	router.HandleFunc("/api/profile/subject", userHandler.LastSubject).Methods(http.MethodGet)
	// *  установка последнего предмета
	router.HandleFunc("/api/profile/subject/{course_id}", userHandler.LastSubject).Methods(http.MethodPost)
	
	// !!! Получение пользователей которую находяться на определенном предмете по id предмета
	router.HandleFunc("/api/profiles/on/subject/{course_id}",userHandler.GetUserOnSubject).Methods(http.MethodGet)
	
	
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

	router.HandleFunc("/api/test/check/{test_id}/{subject_id}", testHandler.CheckQuestion).Methods(
		http.MethodPost,
	)
	router.HandleFunc("/api/test/for/{subject_id}", testHandler.CompletedTestBySubject).Methods(
		http.MethodGet,
	)

	router.HandleFunc("/api/test/result/{test_id}", testHandler.CompletedTest).Methods(
		http.MethodPost,
		http.MethodDelete,
		http.MethodPut,
	)

	// * Общее колчество пройденных тестов

	router.HandleFunc("/api/test/{user_id}/all",testHandler.GetAllCompleted).Methods(
		http.MethodGet,
	)

	// ! ПОЛУЧЕНИЕ СЕРТИФИКАТА
	router.HandleFunc("/api/certificate/{subject_id}", subjectHandler.Certificate).Methods(
		http.MethodGet,
		http.MethodPost,
		http.MethodDelete,
	)

	// ! Эндпоинт для публичного доступа к изображениям)
	// ! НЕ ТРОГАТЬ Я НЕ ЗНАЮ КАК НО ОНО РАБОТАЕТ
	router.HandleFunc("/images", subjectHandler.Image).Methods(http.MethodGet)



	router.HandleFunc("/videos", subjectHandler.Video).Methods(http.MethodGet)


	return router
}

