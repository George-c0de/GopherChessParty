package jwt

import (
	"github.com/golang-jwt/jwt"
	"time"
)

var jwtSecret = []byte("your-secret-key") // секретный ключ лучше хранить в конфигурации

// GenerateToken создает новый JWT-токен для пользователя.
func GenerateToken(username string) (string, error) {
	// Создаем новый токен с методом подписи HS256
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // срок действия токена: 72 часа

	// Подписываем токен секретным ключом
	return token.SignedString(jwtSecret)
}
