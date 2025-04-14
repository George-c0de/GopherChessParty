package interfaces

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/models"
	"github.com/golang-jwt/jwt"
)

type IService interface {
	CreateUser(data *dto.CreateUser) (*models.User, error)
	ValidPassword(data dto.AuthenticateUser) bool
	IsValidateToken(tokenString string) (*jwt.Token, bool)
	GetUsers() ([]*models.User, error)
	GenerateToken(email string) (string, error)
}
