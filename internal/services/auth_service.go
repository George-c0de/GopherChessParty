package services

import (
	"log/slog"
	"time"

	"GopherChessParty/internal/errors"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	GenerateToken(email string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GeneratePassword(rawPassword string) (string, error)
	IsValidPassword(hashedPassword string, plainPassword string) bool
}

type AuthService struct {
	log       *slog.Logger
	jwtSecret string
	exp       time.Duration
}

func (s *AuthService) GenerateToken(email string) (string, error) {
	// Создаем новый токен с методом подписи HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": email,
		"exp":      time.Now().Add(s.exp).Unix(),
	})

	// Подписываем токен секретным ключом
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи соответствует ожиданиям
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrValidateToken
		}
		return []byte(s.jwtSecret), nil
	})
}

func (s *AuthService) IsValidPassword(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func (s *AuthService) GeneratePassword(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("Failed to hash password", err)
		return "", err
	}
	return string(hashedPassword), nil
}
