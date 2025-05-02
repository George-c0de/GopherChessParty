package interfaces

import (
	"GopherChessParty/internal/dto"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type IService interface {
	IUserService
	IGameService
	IAuthService
	IMatchService
	CreateUser(data *dto.CreateUser) (*dto.User, error)
	ValidPassword(data dto.AuthenticateUser) (*uuid.UUID, bool)
	IsValidateToken(tokenString string) (*jwt.Token, bool)
	SearchPlayerConn()
	RegisterUser(data *dto.CreateUser) (*dto.User, error)
	MoveGameStr(gameID uuid.UUID, move string, player *dto.PlayerConn) bool
	SetConnGame(GameID uuid.UUID, player *dto.PlayerConn) error
	GetGameInfoMemory(GameID uuid.UUID, ok bool, move string) (map[string]interface{}, error)
}
