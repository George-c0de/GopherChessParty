package interfaces

import (
	"GopherChessParty/internal/dto"
	"github.com/google/uuid"
)

type IMatchService interface {
	AddUser(player dto.PlayerConn)
	CreatePair() (uuid.UUID, uuid.UUID)
	SearchPlayerConn()
}
