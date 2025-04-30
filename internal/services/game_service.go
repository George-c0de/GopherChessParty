package services

import (
	"GopherChessParty/ent"
	"GopherChessParty/internal/interfaces"
	"github.com/corentings/chess/v2"
	"github.com/google/uuid"
)

type GameService struct {
	log        interfaces.ILogger
	repository interfaces.IGameRepo
}

func (m *GameService) CreateGame(playerID1, playerID2 uuid.UUID) (*ent.Chess, error) {
	return m.repository.Create(playerID1, playerID2)
}

func NewGameService(
	log interfaces.ILogger,
	repository interfaces.IGameRepo,
) interfaces.IGameService {
	return &GameService{
		log:        log,
		repository: repository,
	}
}

func (m *GameService) Move(game *chess.Game, move string) (bool, error) {
	err := game.PushMove(move, nil)
	if err != nil {
		m.log.Error(err)
		return false, err
	}
	return true, nil
}

func (m *GameService) GetGames(userID uuid.UUID) ([]*ent.Chess, error) {
	return m.repository.GetGames(userID)
}
