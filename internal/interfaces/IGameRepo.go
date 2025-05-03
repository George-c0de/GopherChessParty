package interfaces

import (
	"GopherChessParty/ent"
	"GopherChessParty/ent/chess"
	"GopherChessParty/internal/dto"
	"github.com/google/uuid"
)

type IGameRepo interface {
	GetGames(userID uuid.UUID) ([]*dto.GameHistory, error)
	Create(playerID1, playerID2 uuid.UUID) (*ent.Chess, error)
	GetGameById(gameId uuid.UUID) (*dto.Match, error)
	GetStatus(GameID uuid.UUID) chess.Status
	UpdateGameResult(GameId uuid.UUID, status chess.Status, result chess.Result) error
	SaveMove(GameID uuid.UUID, move string, UserID uuid.UUID, numMove int) (*ent.GameHistory, error)
}
