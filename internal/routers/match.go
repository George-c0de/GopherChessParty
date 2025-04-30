package routers

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddWebSocket регистрирует эндпоинт поиска соперника по WebSocket
func AddWebSocket(rg *gin.RouterGroup, service interfaces.IService, log interfaces.ILogger) {
	// Важно: не доверяйте всем прокси (см. https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies)
	rg.Use(middleware.WebSocketTokenMiddleware())
	rg.Use(middleware.JWTAuthMiddleware(service))
	rg.GET("/ws/search", SearchMatchHandler(log))
}

// SearchMatchHandler — обработчик WebSocket для матчмейкинга
func SearchMatchHandler(logger interfaces.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := CreateWebSocket(c)
		if err != nil {
			logger.Error(err)
		}

		// Получение пользователя
		userId, err := middleware.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		player := &dto.PlayerConn{UserID: userId, Conn: conn}

		service := GetService(c)
		service.AddUser(player)
	}
}
