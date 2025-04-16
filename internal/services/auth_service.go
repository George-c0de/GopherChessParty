package services

import (
	"github.com/google/uuid"
	"time"

	"GopherChessParty/internal/errors"
	"GopherChessParty/internal/interfaces"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	log       interfaces.ILogger
	jwtSecret string
	exp       time.Duration
}

func (s *AuthService) GenerateToken(userId uuid.UUID) (string, error) {
	// Создаем новый токен с методом подписи HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(s.exp).Unix(),
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
		s.log.Error(err)
		return "", err
	}
	return string(hashedPassword), nil
}
