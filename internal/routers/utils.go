package routers

import (
	"net/http"
	"time"

	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetService(c *gin.Context) interfaces.IService {
	svc, exists := c.Get("service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service not found"})
		return nil
	}
	// Приводим к нужному типу
	return svc.(*services.Service)
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userId, exists := c.Get("id")
	if !exists {
		return uuid.UUID{}, false
	}
	return userId.(uuid.UUID), true
}

func BindJSON[T any](c *gin.Context) (*T, error) {
	var data T
	if err := c.BindJSON(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func CreateWebSocket(c *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}
	err = conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return nil, err
	}
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	// Тикер для отправки ping
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return // соединение закрыто или ошибка
			}
		}
	}()
	return conn, nil
}
