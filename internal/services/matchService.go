package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Константы для настройки WebSocket соединения
const (
	// Интервал отправки пинг-сообщений для поддержания соединения активным
	pingPeriod = 30 * time.Second
	// Максимальное время ожидания ответа на пинг перед закрытием соединения
	pongWait = 60 * time.Second
)

// MatchService управляет поиском и созданием пар для игры
type MatchService struct {
	log     interfaces.ILogger
	channel chan dto.PlayerConn // Канал для получения новых подключений
	exists  chan struct{}       // Сигнальный канал для оповещения о новых игроках
	queue   []*dto.PlayerConn   // Очередь ожидающих игроков
	queueMu sync.Mutex          // Мьютекс для безопасного доступа к очереди
}

// NewMatchService создает новый сервис матчмейкинга
func NewMatchService(log interfaces.ILogger) *MatchService {
	return &MatchService{
		log:     log,
		channel: make(chan dto.PlayerConn, 100), // Буферизованный канал на 100 подключений
		exists:  make(chan struct{}, 1),         // Буферизованный канал для сигналов
	}
}

func (m *MatchService) CheckPair() bool {
	m.queueMu.Lock()
	defer m.queueMu.Unlock()
	return len(m.queue) == 2
}

// GetExistsChannel возвращает канал для оповещения о новых игроках
func (m *MatchService) GetExistsChannel() <-chan struct{} {
	return m.exists
}

// setupPing настраивает механизм поддержания WebSocket соединения активным
func (m *MatchService) setupPing(conn *websocket.Conn) {
	// Обработчик входящих пинг-сообщений
	conn.SetPingHandler(func(string) error {
		// Отправляем понг в ответ на пинг
		return conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(pongWait))
	})

	// Обработчик входящих понг-сообщений
	conn.SetPongHandler(func(string) error {
		// Обновляем таймаут чтения при получении понга
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	// Запускаем горутину для периодической отправки пингов
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		// Каждые pingPeriod отправляем пинг
		for range ticker.C {
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(pongWait)); err != nil {
				m.log.Error(fmt.Errorf("ошибка отправки пинга: %w", err))
				return
			}
		}
	}()
}

// AddUser добавляет нового игрока в очередь ожидания
func (m *MatchService) AddUser(player *dto.PlayerConn) {
	m.queueMu.Lock()
	defer m.queueMu.Unlock()
	m.queue = append(m.queue, player)
	m.setupPing(player.Conn) // Настраиваем пинг для нового соединения
	m.exists <- struct{}{}   // Сигнализируем о новом игроке
}

// sendMessage отправляет сообщение игроку через WebSocket
func (m *MatchService) sendMessage(player *dto.PlayerConn, message []byte) error {
	err := player.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		m.log.Error(err)
		return err
	}
	return nil
}

func (m *MatchService) ReturnPlayerID() (*dto.PlayerConn, *dto.PlayerConn, error) {
	m.queueMu.Lock()
	defer m.queueMu.Unlock()
	player1, player2 := m.queue[0], m.queue[1]
	m.queue = m.queue[2:]
	return player1, player2, nil
}

func (m *MatchService) SendGemID(player *dto.PlayerConn, gameID uuid.UUID) error {
	message := make(map[string]uuid.UUID)
	message["gameID"] = gameID
	messageByte, err := json.Marshal(message)
	if err != nil {
		m.log.Error(err)
		return err
	}
	err = m.sendMessage(player, messageByte)
	if err != nil {
		return err
	}
	return nil
}

func (m *MatchService) CloseConnection(player *dto.PlayerConn) error {
	err := player.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (m *MatchService) SendMove(player *dto.PlayerConn, move string) error {
	answer := map[string]interface{}{
		"move": move,
	}
	response, err := json.Marshal(answer)
	if err != nil {
		m.log.Error(err)
		return err
	}
	errSend := player.Conn.WriteMessage(websocket.TextMessage, response)
	if errSend != nil {
		m.log.Error(err)
		return err
	}
	return nil
}
