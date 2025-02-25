package main

import (
	"fmt"
	"github.com/George-c0de/GopherChessParty/internal/config"
	"github.com/George-c0de/GopherChessParty/internal/logger"
	chess "github.com/George-c0de/GopherChessParty/internal/services"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Info("starting server")
	manager := chess.New(log)
	myGame := manager.CreateGame()
	manager.Move(myGame, "e4")
	manager.Move(myGame, "e5")
	manager.Move(myGame, "Nf3")
	fmt.Printf(myGame.Position().Board().Draw())
}
