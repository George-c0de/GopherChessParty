package main

import (
	"fmt"

	"GopherChessParty/internal/config"
	"GopherChessParty/internal/logger"
	"GopherChessParty/internal/repository"
	"GopherChessParty/internal/routers"
	"GopherChessParty/internal/services"
)

func main() {
	// Загрузка конфига
	cfg := config.MustLoad()

	// Создание Логгера
	log := logger.New(cfg.Application)

	// Подключение к базе данных и создание репозитория
	connection := repository.MustNewConnection(cfg.Database)

	// Создание репозиториев
	userRepo := repository.NewUserRepository(log, connection)
	gameRepo := repository.NewGameRepository(log, connection)

	// Создание сервиса
	userService := services.NewUserService(log, userRepo)
	gameService := services.NewGameService(log, gameRepo)
	authService := services.NewAuthService(log, cfg.Auth)
	matchService := services.NewMatchService(log)

	service := services.NewService(userService, gameService, authService, matchService, log)

	// Создание экземпляра Gin
	router := routers.New(service, log)

	err := router.Run(fmt.Sprintf(":%d", cfg.Application.Port))
	if err != nil {
		panic(err)
	}
}
