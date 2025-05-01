package interfaces

import (
	"GopherChessParty/ent"
	"GopherChessParty/internal/dto"
	"github.com/google/uuid"
)

type IGameService interface {
	GetGames(userID uuid.UUID) ([]*ent.Chess, error)
	CreateGame(playerID1, playerID2 uuid.UUID) (*ent.Chess, error)
	GetGameByID(gameID uuid.UUID) (*dto.Match, error)
	MoveGame(GameID uuid.UUID, move string, player *dto.PlayerConn) error
	SetPlayer(GameID uuid.UUID, player *dto.PlayerConn) error
	IsConnectPlayers(GameID uuid.UUID) bool
	GetOpponent(gameID uuid.UUID) *dto.PlayerConn
}
