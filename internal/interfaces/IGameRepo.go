package interfaces

import (
	"GopherChessParty/ent"
	"GopherChessParty/ent/chess"
	"GopherChessParty/internal/dto"
	"github.com/google/uuid"
)

type IGameRepo interface {
	Games(userID uuid.UUID) ([]*dto.GameHistory, error)
	Create(playerID1, playerID2 uuid.UUID) (*ent.Chess, error)
	GameById(gameId uuid.UUID) (*dto.Match, error)
	Status(GameID uuid.UUID) chess.Status
	UpdateGame(GameId uuid.UUID, status chess.Status, result chess.Result) error
	SaveMove(GameID uuid.UUID, move string, UserID uuid.UUID, numMove int) (*ent.GameHistory, error)
}
