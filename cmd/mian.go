package main

import (
	"context"
	"ggkit_learn_service/internals/app"
	"ggkit_learn_service/internals/cfg"
	"log"
	"os"
	"os/signal"
)

func main() {
	config := cfg.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	server := app.NewServer(config, ctx)

	go func() { //горутина для ловли сообщений системы
		oscall := <-c //если таки что то пришло
		log.Printf("system call:%#+v", oscall)
		server.Shutdown() //выключаем сервер
		cancel()          //отменяем контекст
	}()
	server.Serve() //запускаем сервер

}
