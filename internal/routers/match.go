package routers

import (
	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type PlayerConn struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
}

var (
	queue   []*PlayerConn
	queueMu sync.Mutex
)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	pongWait   = 60 * time.Second
	pingPeriod = 30 * time.Second
)

// AddWebSocket регистрирует эндпоинт поиска соперника по WebSocket
func AddWebSocket(rg *gin.RouterGroup, service interfaces.IService, log interfaces.ILogger) {
	// Важно: не доверяйте всем прокси (см. https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies)
	rg.Use(middleware.JWTAuthMiddleware(service))
	rg.GET("/ws/search", SearchMatchHandler(log))
}

// SearchMatchHandler — обработчик WebSocket для матчмейкинга
func SearchMatchHandler(logger interfaces.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error(err)
			return
		}
		defer conn.Close()

		userId, err := middleware.GetUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		// Канал и once для сигнализации о закрытии
		done := make(chan struct{})
		var closeOnce sync.Once

		// Отправка ping
		go keepAlivePing(conn, pingPeriod, done)

		player := &PlayerConn{UserID: userId, Conn: conn}

		defer dequeueByID(player.UserID)

		if opponent, found := enqueue(player); found {
			logger.Info("Матч найден: %s vs %s", player.UserID, opponent.UserID)
			message := gin.H{
				"type": "match_found",
				"data": gin.H{"opponent": opponent.UserID},
			}
			notifyBoth(player, opponent, message)
			dequeueByID(opponent.UserID)
			opponent.Conn.Close()
			closeOnce.Do(func() { close(done) })
			return
		}

		// Если соперник ещё не найден — ждём сигнала о завершении
		<-done
	}
}

// enqueue — добавляет игрока в очередь или возвращает готового соперника
func enqueue(p *PlayerConn) (*PlayerConn, bool) {
	queueMu.Lock()
	defer queueMu.Unlock()

	if len(queue) == 0 {
		queue = append(queue, p)
		return nil, false
	}

	opponent := queue[0]
	queue = queue[1:]
	return opponent, true
}

// dequeueByID — удаляет игрока из очереди по ID
func dequeueByID(id uuid.UUID) {
	queueMu.Lock()
	defer queueMu.Unlock()

	for i, p := range queue {
		if p.UserID == id {
			queue = append(queue[:i], queue[i+1:]...)
			return
		}
	}
}

// keepAlivePing регулярно шлёт Ping и завершает по сигналу done
func keepAlivePing(conn *websocket.Conn, interval time.Duration, done chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Пинг не отправлен, соединение закрыто:", err)
				return
			}
		case <-done:
			return
		}
	}
}

// notifyBoth — отправляет JSON-сообщение обоим игрокам
func notifyBoth(p1, p2 *PlayerConn, msg any) {
	res, _ := json.Marshal(msg)
	if err := p1.Conn.WriteMessage(websocket.TextMessage, res); err != nil {
		log.Println("Ошибка при отправке p1:", err)
	}
	if err := p2.Conn.WriteMessage(websocket.TextMessage, res); err != nil {
		log.Println("Ошибка при отправке p2:", err)
	}
}
