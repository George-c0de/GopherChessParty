package middleware

import (
	"net/http"
	"strings"

	"GopherChessParty/internal/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// JWTAuthMiddleware проверяет валидность JWT-токена.
func JWTAuthMiddleware(service interfaces.IService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		// Ожидается формат "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Парсим токен
		token, ok := service.IsValidateToken(tokenString)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Можно извлечь claims и положить в контекст запроса, если потребуется
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("id", claims["id"])
		}

		c.Next()
	}
}
