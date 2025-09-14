package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	srv := server.NewServer(logger)

	logger.Println("Сервер запущен")

	err := srv.Server.ListenAndServe()
	if err != nil {
		logger.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
