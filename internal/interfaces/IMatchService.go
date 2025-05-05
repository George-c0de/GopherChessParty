package interfaces

import (
	"GopherChessParty/internal/dto"
	"github.com/google/uuid"
)

type IMatchService interface {
	CheckPair() bool
	ExistsChannel() <-chan struct{}
	AddUser(player *dto.PlayerConn) error
	ReturnPlayers() (player1, player2 *dto.PlayerConn)
	SendGemID(player *dto.PlayerConn, gameID uuid.UUID) error
	CloseConnection(player *dto.PlayerConn) error
	SendMove(player *dto.PlayerConn, move string) error
	SendMessage(player *dto.PlayerConn, message map[string]interface{}) error
	ExitPlayerAdd(playerID uuid.UUID)
}
