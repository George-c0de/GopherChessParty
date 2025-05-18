package services

import (
	"time"

	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/errors"
	"GopherChessParty/internal/interfaces"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	log interfaces.ILogger
	dto.AuthConfig
}

func NewAuthService(log interfaces.ILogger, cfg dto.AuthConfig) interfaces.IAuthService {
	return &AuthService{
		log,
		cfg,
	}
}

func (s *AuthService) GetUserIDFromRefresh(refreshToken string) (uuid.UUID, error) {
	token, err := s.ValidateRefresh(refreshToken)
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, errors.ErrValidateToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.ErrValidateToken
	}

	idRaw, ok := claims["id"].(string)
	if !ok {
		return uuid.Nil, errors.ErrValidateToken
	}
	userID, err := uuid.Parse(idRaw)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (s *AuthService) generateAccess(userId uuid.UUID) (string, error) {
	// Создаем новый токен с методом подписи HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(s.ExpAccess).Unix(),
	})

	// Подписываем токен секретным ключом
	return token.SignedString([]byte(s.SecretAccess))
}

func (s *AuthService) generateRefresh(userId uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(s.ExpRefresh).Unix(),
	})

	// Подписываем токен секретным ключом
	return token.SignedString([]byte(s.SecretRefresh))
}

func (s *AuthService) GenerateTokenPair(userId uuid.UUID) (*dto.TokenPair, error) {
	accessToken, err := s.generateAccess(userId)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.generateRefresh(userId)
	if err != nil {
		return nil, err
	}
	return &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *AuthService) RefreshAccessToken(refreshToken string) (*dto.TokenPair, error) {
	userID, err := s.GetUserIDFromRefresh(refreshToken)
	if err != nil {
		return nil, err
	}
	return s.GenerateTokenPair(userID)
}

func (s *AuthService) ValidateAccess(access string) (*jwt.Token, error) {
	return jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи соответствует ожиданиям
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.log.Error(errors.ErrValidateToken)
			return nil, errors.ErrValidateToken
		}
		return []byte(s.SecretAccess), nil
	})
}
func (s *AuthService) ValidateRefresh(refresh string) (*jwt.Token, error) {
	return jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи соответствует ожиданиям
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.log.Error(errors.ErrValidateToken)
			return nil, errors.ErrValidateToken
		}
		return []byte(s.SecretRefresh), nil
	})
}

func (s *AuthService) IsValidPassword(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func (s *AuthService) GeneratePassword(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error(err)
		return "", err
	}
	return string(hashedPassword), nil
}
