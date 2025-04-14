package middleware

import (
	"GopherChessParty/internal/interfaces"
	"github.com/gin-gonic/gin"
)

func ServiceMiddleware(service interfaces.IService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Добавляем сервис в контекст с ключом "service"
		c.Set("service", service)
		c.Next()
	}
}
