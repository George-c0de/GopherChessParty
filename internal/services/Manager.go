package chess

import (
	"github.com/corentings/chess/v2"
	"github.com/google/uuid"
	"log/slog"
)

type Manager struct {
	log      *slog.Logger
	nowGames map[uuid.UUID]*chess.Game
}

func New(log *slog.Logger) *Manager {
	return &Manager{log, make(map[uuid.UUID]*chess.Game)}
}

func (m *Manager) CreateGame() *chess.Game {
	random, err := uuid.NewRandom()
	if err != nil {
		m.log.Error(err.Error())
		panic(err)
	}
	m.nowGames[random] = chess.NewGame()
	return m.nowGames[random]
}

func (m *Manager) Move(game *chess.Game, move string) (error, bool) {
	err := game.PushMove(move, nil)
	if err != nil {
		m.log.Error(err.Error())
		return err, false
	}
	return nil, true
}
