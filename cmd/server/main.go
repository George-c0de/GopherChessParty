package main

import (
	"fmt"
	"github.com/George-c0de/GopherChessParty/internal/config"
	"github.com/George-c0de/GopherChessParty/internal/dto"
	"github.com/George-c0de/GopherChessParty/internal/logger"
	chess "github.com/George-c0de/GopherChessParty/internal/services"
	"github.com/George-c0de/GopherChessParty/internal/storage"
	"github.com/George-c0de/GopherChessParty/internal/storage/postgres"
	"math/rand"
	"strconv"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Info("starting server")
	manager := chess.New(log)
	database, _ := postgres.New(log)
	repository := storage.NewRepository(database)
	users := repository.GetUsers()

	for i := 0; i < len(users); i++ {
		user := users[i]
		fmt.Println(user)
	}
	createUser, _ := repository.CreateUser(&dto.CreateUser{
		Name:     "Test" + strconv.Itoa(rand.Intn(100)),
		Email:    fmt.Sprintf("test%d@example.com", rand.Intn(1000)),
		Password: "test",
	})
	fmt.Println(createUser)
	myGame := manager.CreateGame()
	manager.Move(myGame, "e4")
	manager.Move(myGame, "e5")
	manager.Move(myGame, "Nf3")
	fmt.Printf(myGame.Position().Board().Draw())
}
