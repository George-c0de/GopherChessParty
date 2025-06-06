package routers

import (
	"net/http"
	"time"

	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 5 * time.Minute
	pongWait       = 1 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
	conn.SetReadLimit(maxMessageSize)
	err = conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return nil, err
	}
	conn.SetPongHandler(func(string) error {
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}
		return nil
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
