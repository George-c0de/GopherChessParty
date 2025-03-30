package services

import (
	"log/slog"

	"github.com/corentings/chess/v2"
	"github.com/google/uuid"
)

type IGameService interface {
	CreateGame() *chess.Game
	Move(game *chess.Game, move string) (error, bool)
}
type GameService struct {
	log      *slog.Logger
	nowGames map[uuid.UUID]*chess.Game
}

func (m *GameService) CreateGame() *chess.Game {
	random, err := uuid.NewRandom()
	if err != nil {
		m.log.Error(err.Error())
		panic(err)
	}

	m.nowGames[random] = chess.NewGame()
	return m.nowGames[random]
}

func (m *GameService) Move(game *chess.Game, move string) (error, bool) {
	err := game.PushMove(move, nil)
	if err != nil {
		m.log.Error(err.Error())
		return err, false
	}
	return nil, true
}
