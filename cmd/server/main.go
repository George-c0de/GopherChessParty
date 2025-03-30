package main

import (
	"GopherChessParty/internal/config"
	"GopherChessParty/internal/logger"
	"GopherChessParty/internal/routers"
	"GopherChessParty/internal/services"
	"GopherChessParty/internal/storage/ent_repository"
)

func main() {
	// Загрузка конфига
	cfg := config.MustLoad()

	// Создание Логгера
	log := logger.SetupLogger(cfg.Env)

	// Подключение к базе данных и создание репозитория
	repository := ent_repository.MustNewRepository(cfg.Database, log)

	// Создание сервиса
	service := services.New(log, repository, cfg.Auth)

	// Создание экземпляра Gin
	router := routers.GetRoutes(service)

	router.Run(":8000")

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
