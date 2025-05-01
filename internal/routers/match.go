package routers

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/middleware"
	"encoding/json"
	"net/http"

	chesslib "github.com/corentings/chess/v2"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
)

// AddWebSocket регистрирует эндпоинт поиска соперника по WebSocket
func AddWebSocket(rg *gin.RouterGroup, service interfaces.IService, log interfaces.ILogger) {
	// Важно: не доверяйте всем прокси (см. https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies)
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

func MoveGame(logger interfaces.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := CreateWebSocket(c)
		if err != nil {
			logger.Error(err)
			return
		}
		defer conn.Close()

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
		game := chesslib.NewGame()
		service.SetConnGame(gameID, player)
		// Запускаем горутину для чтения сообщений
		go func() {
			for {
				messageType, message, err := conn.ReadMessage()
				if err != nil {
					logger.Error(err)
					return
				}

				if messageType == websocket.TextMessage {
					// Обрабатываем ход
					move := string(message)
					status, ok := service.MoveGameStr(gameID, move, player)

					if ok {
						// Выводим текущее состояние доски
						logger.Info("\n" + game.CurrentPosition().Board().Draw())
					}

					// Отправляем ответ
					answer := map[string]interface{}{
						"ok":     ok,
						"status": status,
					}
					if !ok {
						answer["message"] = "Недопустимый ход"
					}

					response, err := json.Marshal(answer)
					if err != nil {
						logger.Error(err)
						continue
					}

					if err := conn.WriteMessage(websocket.TextMessage, response); err != nil {
						logger.Error(err)
						return
					}
				}
			}
		}()

		// Держим соединение открытым
		select {}
	}
}
