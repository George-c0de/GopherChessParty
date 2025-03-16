package main

import (
	"fmt"
	"github.com/George-c0de/GopherChessParty/internal/config"
	"github.com/George-c0de/GopherChessParty/internal/logger"
	chess "github.com/George-c0de/GopherChessParty/internal/services"
	"github.com/George-c0de/GopherChessParty/internal/storage"
	"github.com/George-c0de/GopherChessParty/internal/storage/postgres"
)

func main() {
	// Загрузка конфига
	cfg := config.MustLoad()

	// Создание Логгера
	log := logger.SetupLogger(cfg.Env)

	// Создание сервисов
	manager := chess.New(log)

	// Подключение к базе данных и создание репозитория
	database := postgres.MustNew(log, cfg.Database)
	repository := storage.NewRepository(database)
	users := repository.GetUsers()

	for i := 0; i < len(users); i++ {
		user := users[i]
		fmt.Println(user)
	}
	myGame := manager.CreateGame()
	manager.Move(myGame, "e4")
	manager.Move(myGame, "e5")
	manager.Move(myGame, "Nf3")
	fmt.Printf(myGame.Position().Board().Draw())
}
