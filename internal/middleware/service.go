package middleware

import (
	"GopherChessParty/internal/services"
	"github.com/gin-gonic/gin"
)

func ServiceMiddleware(service services.IService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Добавляем сервис в контекст с ключом "service"
		c.Set("service", service)
		c.Next()
	}
}
