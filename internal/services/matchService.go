package services

import (
	"encoding/json"
	"sync"

	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

// AddUser добавляет нового игрока в очередь ожидания
func (m *MatchService) AddUser(player *dto.PlayerConn) {
	m.queueMu.Lock()
	defer m.queueMu.Unlock()
	m.queue = append(m.queue, player)
	m.exists <- struct{}{} // Сигнализируем о новом игроке
}

func (m *MatchService) ReturnPlayerID() (*dto.PlayerConn, *dto.PlayerConn) {
	m.queueMu.Lock()
	defer m.queueMu.Unlock()
	player1, player2 := m.queue[0], m.queue[1]
	m.queue = m.queue[2:]
	return player1, player2
}

func (m *MatchService) SendGemID(player *dto.PlayerConn, gameID uuid.UUID) error {
	message := map[string]interface{}{
		"gameID": gameID,
	}
	err := m.SendMessage(player, message)
	if err != nil {
		m.log.Error(err)
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
	err := m.SendMessage(player, answer)
	if err != nil {
		return err
	}
	return nil
}

func (m *MatchService) SendMessage(player *dto.PlayerConn, message map[string]interface{}) error {
	response, err := json.Marshal(message)
	if err != nil {
		m.log.Error(err)
		return err
	}
	errSend := player.Conn.WriteMessage(websocket.TextMessage, response)
	if errSend != nil {
		m.log.Error(errSend)
		return errSend
	}
	return nil
}
