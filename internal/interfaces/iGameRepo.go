package interfaces

import (
	"GopherChessParty/ent"
	"github.com/google/uuid"
)

type IGameRepo interface {
	GetGames(userID uuid.UUID) ([]*ent.Chess, error)
}
