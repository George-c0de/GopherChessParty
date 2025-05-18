package interfaces

import (
	"GopherChessParty/internal/dto"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type IAuthService interface {
	GenerateTokenPair(userId uuid.UUID) (*dto.TokenPair, error)
	RefreshAccessToken(refreshToken string) (*dto.TokenPair, error)
	ValidateAccess(access string) (*jwt.Token, error)
	ValidateRefresh(refresh string) (*jwt.Token, error)
	GeneratePassword(rawPassword string) (string, error)
	IsValidPassword(hashedPassword string, plainPassword string) bool
}
