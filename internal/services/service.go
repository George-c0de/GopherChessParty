package services

import (
	"github.com/George-c0de/GopherChessParty/internal/models"
	"github.com/George-c0de/GopherChessParty/internal/storage"
	"github.com/corentings/chess/v2"
	"github.com/google/uuid"
	"log/slog"
)

type Service struct {
	log        *slog.Logger
	nowGames   map[uuid.UUID]*chess.Game
	Repository *storage.Repository
}

func New(log *slog.Logger, repository *storage.Repository) *Service {
	return &Service{log, make(map[uuid.UUID]*chess.Game), repository}
}

func (m *Service) CreateGame() *chess.Game {
	random, err := uuid.NewRandom()
	if err != nil {
		m.log.Error(err.Error())
		panic(err)
	}
	m.nowGames[random] = chess.NewGame()
	return m.nowGames[random]
}

func (m *Service) GetUsers() []*models.User {
	return m.Repository.GetUsers()
}
func (m *Service) Move(game *chess.Game, move string) (error, bool) {
	err := game.PushMove(move, nil)
	if err != nil {
		m.log.Error(err.Error())
		return err, false
	}
	return nil, true
}
