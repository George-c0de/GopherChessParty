package main

import (
	"fmt"
	"github.com/George-c0de/GopherChessParty/internal/config"
	"github.com/George-c0de/GopherChessParty/internal/logger"
	"github.com/George-c0de/GopherChessParty/internal/services"
	"github.com/George-c0de/GopherChessParty/internal/storage"
	"github.com/George-c0de/GopherChessParty/internal/storage/postgres"
)

func main() {
	// Загрузка конфига
	cfg := config.MustLoad()

	// Создание Логгера
	log := logger.SetupLogger(cfg.Env)

	// Подключение к базе данных и создание репозитория
	database := postgres.MustNew(log, cfg.Database)
	repository := storage.NewRepository(database)

	// Создание сервиса
	service := services.New(log, repository)
	users := service.GetUsers()
	for i := 0; i < len(users); i++ {
		user := users[i]
		fmt.Println(user)
	}
	myGame := service.CreateGame()
	service.Move(myGame, "e4")
	service.Move(myGame, "e5")
	service.Move(myGame, "Nf3")
	fmt.Printf(myGame.Position().Board().Draw())
}
