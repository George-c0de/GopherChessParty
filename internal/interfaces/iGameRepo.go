package interfaces

import (
	"GopherChessParty/ent"
	"github.com/google/uuid"
)

type IGameRepo interface {
	GetGames(userID uuid.UUID) ([]*ent.Chess, error)
	Create(playerID1, playerID2 uuid.UUID) (*ent.Chess, error)
}
