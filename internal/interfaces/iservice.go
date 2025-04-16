package interfaces

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type IService interface {
	IUserService
	IGameService
	IAuthService
	CreateUser(data *dto.CreateUser) (*models.User, error)
	ValidPassword(data dto.AuthenticateUser) (*uuid.UUID, bool)
	IsValidateToken(tokenString string) (*jwt.Token, bool)
	GetGames(userId uuid.UUID)
}
