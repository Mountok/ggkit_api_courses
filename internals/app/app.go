package app

import (
	"context"
	"ggkit_learn_service/api"
	"ggkit_learn_service/internals/app/db"
	"ggkit_learn_service/internals/app/handler"
	"ggkit_learn_service/internals/app/processor"
	"ggkit_learn_service/internals/cfg"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AppServer struct {
	config cfg.Cfg
	ctx    context.Context
	serv   *http.Server
	db     *pgxpool.Pool
}

func NewServer(config cfg.Cfg, ctx context.Context) *AppServer { //задаем поля нашего сервера, для его старта нам нужен контекст и конфигурация
	server := new(AppServer)
	server.ctx = ctx
	server.config = config
	return server
}

func (server *AppServer) Serve() {
	log.Println("Server starting")
	log.Println(server.config.GetDBConnetcUrl())
	var err error
	//создаем пул соединений с БД и сохраним его для закрытия при остановке приложения
	server.db, err = pgxpool.Connect(server.ctx, server.config.GetDBConnetcUrl())
	if err != nil {
		log.Fatalln(err)
	}

	subjectStorage := db.NewSubjectStorage(server.db)
	subjectProcessor := processor.NewSubjectProcessor(subjectStorage)
	subjectHandler := handler.NewSubjectHandler(subjectProcessor)

	themeStorage := db.NewThemesStorage(server.db)
	themeProcessor := processor.NewThemesProcessor(themeStorage)
	themeHandler := handler.NewThemesHandler(themeProcessor)

	lessonsStorage := db.NewLessonsStorage(server.db)
	lessonsProcess := processor.NewLessonsProcessor(lessonsStorage)
	lessonsHandler := handler.NewLessonsHanler(lessonsProcess)

	loginStorage := db.NewLoginStorage(server.db)
	loginProcess := processor.NewLoginProcessor(loginStorage)
	loginHandler := handler.NewLoginhandler(loginProcess)

	userStorage := db.NewUserStorage(server.db)
	userProcess := processor.NewUserProcessor(userStorage)
	userHandler := handler.NewUserHandler(userProcess)

	testStorage := db.NewSubjectTestStorage(server.db)
	testProcess := processor.NewSubjectTestProcessor(testStorage)
	testHandler := handler.NewSubjectTestHandler(testProcess)

	commentStorage := db.NewCommentStorage(server.db)
	commentProcess := processor.NewCommentProcessor(commentStorage)
	commentHandler := handler.NewCommentHandler(commentProcess)

	routes := api.CreateRoute(
		subjectHandler,
		themeHandler,
		lessonsHandler,
		loginHandler,
		userHandler,
		testHandler,
		commentHandler,
	)

	server.serv = &http.Server{
		Addr: "0.0.0.0:" + server.config.Port,
		Handler: handlers.CORS(
			handlers.AllowedOrigins([]string{
				"*",
			}),
			handlers.AllowedMethods([]string{
				"GET",
				"POST",
				"PUT",
				"DELETE",
			}),
			handlers.AllowedHeaders([]string{
				"X-Requested-With",
				"Content-Type",
				"Authorization",
			}),
		)(routes),
	}

	log.Println("server started")

	err = server.serv.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}

	return
}

func (server *AppServer) Shutdown() {
	log.Println("server stopped")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	server.db.Close() //закрываем соединение с БД
	defer func() {
		cancel()
	}()
	var err error
	if err = server.serv.Shutdown(ctxShutDown); err != nil { //выключаем сервер, с ограниченным по времени контекстом
		log.Fatalf("server Shutdown Failed:%s", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
}
