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
		panic(err)
	}
	m.nowGames[random] = chess.NewGame()
	return m.nowGames[random]
}
