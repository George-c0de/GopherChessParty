package interfaces

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type IAuthService interface {
	GenerateToken(userId uuid.UUID) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GeneratePassword(rawPassword string) (string, error)
	IsValidPassword(hashedPassword string, plainPassword string) bool
}
