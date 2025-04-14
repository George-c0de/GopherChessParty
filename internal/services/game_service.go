package services

import (
	"GopherChessParty/internal/interfaces"
	"github.com/corentings/chess/v2"
	"github.com/google/uuid"
)

type GameService struct {
	log      interfaces.ILogger
	nowGames map[uuid.UUID]*chess.Game
}

func (m *GameService) CreateGame() *chess.Game {
	random, err := uuid.NewRandom()
	if err != nil {
		m.log.Error(err)
		panic(err)
	}

	m.nowGames[random] = chess.NewGame()
	return m.nowGames[random]
}

func (m *GameService) Move(game *chess.Game, move string) (bool, error) {
	err := game.PushMove(move, nil)
	if err != nil {
		m.log.Error(err)
		return false, err
	}
	return true, nil
}
