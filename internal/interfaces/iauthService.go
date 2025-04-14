package interfaces

import "github.com/golang-jwt/jwt"

type IAuthService interface {
	GenerateToken(email string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GeneratePassword(rawPassword string) (string, error)
	IsValidPassword(hashedPassword string, plainPassword string) bool
}
