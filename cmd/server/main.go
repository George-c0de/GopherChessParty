package main

import (
	"GopherChessParty/internal/config"
	"GopherChessParty/internal/logger"
	"GopherChessParty/internal/routers"
	"GopherChessParty/internal/services"
	"GopherChessParty/internal/storage"
)

func main() {
	// Загрузка конфига
	cfg := config.MustLoad()

	// Создание Логгера
	log := logger.NewLogger(cfg.Env)

	// Подключение к базе данных и создание репозитория
	repository := storage.MustNewRepository(cfg.Database)

	// Создание репозиториев
	userRepo := storage.NewUserRepository(log, repository)
	gameRepo := storage.NewGameRepository(log, repository)

	// Создание сервиса
	userService := services.NewUserService(log, userRepo)
	gameService := services.NewGameService(log, gameRepo)
	authService := services.NewAuthService(log, cfg.Auth)
	matchService := services.NewMatchService(log)

	service := services.NewService(userService, gameService, authService, matchService, log)

	// Создание экземпляра Gin
	router := routers.GetRoutes(service, log)

	err := router.Run(":8000")
	if err != nil {
		panic(err)
	}
	//users := service.GetUsers()
	//
	//for i := 0; i < len(users); i++ {
	//	user := users[i]
	//	fmt.Println(user)
	//}
	//myGame := service.CreateGame()
	//service.Move(myGame, "e4")
	//service.Move(myGame, "e5")
	//service.Move(myGame, "Nf3")
	//fmt.Printf(myGame.Position().Board().Draw())
}
