package routers

import (
	"net/http"

	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// AddWebSocket регистрирует эндпоинт поиска соперника по WebSocket
func AddWebSocket(rg *gin.RouterGroup, service interfaces.IService, log interfaces.ILogger) {
	// Важно: https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	rg.Use(middleware.WebSocketTokenMiddleware())
	rg.Use(middleware.JWTAuthMiddleware(service))
	rg.GET("/ws/search", SearchMatchHandler(log))
	rg.GET("/ws/game/:game_id", MoveGame(log))
}

// SearchMatchHandler — обработчик WebSocket для матчмейкинга
func SearchMatchHandler(logger interfaces.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := CreateWebSocket(c)
		if err != nil {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Получение пользователя
		userId, err := middleware.GetUserID(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		player := &dto.PlayerConn{UserID: userId, Conn: conn}

		service := GetService(c)
		err = service.AddUser(player)
		if err != nil {
			_ = conn.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		player.Conn.SetCloseHandler(func(code int, text string) error {
			service.ExitPlayerAdd(player.UserID)
			return nil
		})
	}
}

func MoveGame(logger interfaces.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := CreateWebSocket(c)
		if err != nil {
			logger.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		gameID, err := uuid.Parse(c.Param("game_id"))
		if err != nil {
			logger.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Получение пользователя
		userID, err := middleware.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		player := &dto.PlayerConn{UserID: userID, Conn: conn}

		service := GetService(c)
		errConn := service.SetConnGame(gameID, player)
		if errConn != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// Запускаем горутину для чтения сообщений
		go func() {
			defer func() {
				err := conn.Close()
				if err != nil {
					logger.Error(err)
					return
				}
			}()

			for {
				messageType, message, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(
						err,
						websocket.CloseGoingAway,
						websocket.CloseAbnormalClosure,
					) {
						logger.Error(err)
					}
					return
				}

				switch messageType {
				case websocket.TextMessage:
					// Обрабатываем ход
					move := string(message)
					ok := service.MoveGameStr(gameID, move, player)
					response, _ := service.GetGameInfoMemory(gameID, ok, "")
					_ = service.SendMessage(player, response)
				case websocket.PingMessage:
					err := conn.WriteMessage(websocket.PongMessage, message)
					if err != nil {
						logger.Error(err)
						return
					}
				}
			}
		}()

		// Отправляем начальное состояние игры
		response, err := service.GetGameInfoMemory(gameID, true, "")
		if err != nil {
			return
		}
		errSend := service.SendMessage(player, response)
		if errSend != nil {
			return
		}

		// Держим соединение открытым
		select {}
	}
}
