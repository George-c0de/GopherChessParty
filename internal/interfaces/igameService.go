package interfaces

import (
	"GopherChessParty/ent"
	"github.com/corentings/chess/v2"
	"github.com/google/uuid"
)

type IGameService interface {
	Move(game *chess.Game, move string) (bool, error)
	GetGames(userID uuid.UUID) ([]*ent.Chess, error)
	CreateGame(playerID1, playerID2 uuid.UUID) (*ent.Chess, error)
}
