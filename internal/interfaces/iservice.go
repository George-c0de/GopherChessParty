package interfaces

import (
	"GopherChessParty/ent"
	"GopherChessParty/ent/chess"
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type IService interface {
	IUserService
	IGameService
	IAuthService
	IMatchService
	CreateUser(data *dto.CreateUser) (*models.User, error)
	ValidPassword(data dto.AuthenticateUser) (*uuid.UUID, bool)
	IsValidateToken(tokenString string) (*jwt.Token, bool)
	GetGamesByUserID(userId uuid.UUID) ([]*ent.Chess, error)
	SearchPlayerConn()
	CreateMatch(playerID1, playerID2 uuid.UUID) (*ent.Chess, error)
	RegisterUser(data *dto.CreateUser) (*models.User, error)
	MoveGameStr(gameID uuid.UUID, move string, player *dto.PlayerConn) error
	SetConnGame(GameID uuid.UUID, player *dto.PlayerConn)
	GetStatusMemory(GameID uuid.UUID) (chess.Status, error)
	GetHistoryMove(id uuid.UUID) []string
	GetGameInfoMemory(GameID uuid.UUID, ok bool) (map[string]interface{}, error)
	MoveValid(GameID uuid.UUID, move string) error
}
