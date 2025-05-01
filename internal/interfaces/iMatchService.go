package interfaces

import (
	"GopherChessParty/internal/dto"
	"github.com/google/uuid"
)

type IMatchService interface {
	CheckPair() bool
	GetExistsChannel() <-chan struct{}
	AddUser(player *dto.PlayerConn)
	ReturnPlayerID() (*dto.PlayerConn, *dto.PlayerConn, error)
	SendGemID(player *dto.PlayerConn, gameID uuid.UUID) error
	CloseConnection(player *dto.PlayerConn) error
	SendMove(player *dto.PlayerConn, move string) error
}
