package services

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type MatchService struct {
	log     interfaces.ILogger
	channel chan dto.PlayerConn
	exists  chan struct{}
	queue   []*dto.PlayerConn
	queueMu sync.Mutex
}

func NewMatchService(log interfaces.ILogger) *MatchService {
	return &MatchService{
		log:     log,
		channel: make(chan dto.PlayerConn, 100),
	}
}

func (m *MatchService) AddUser(player dto.PlayerConn) {
	m.queueMu.Lock()
	defer m.queueMu.Unlock()
	m.queue = append(m.queue, &player)
	m.exists <- struct{}{}
}

func (m *MatchService) CreatePair() (uuid.UUID, uuid.UUID) {
	m.queueMu.Lock()
	defer m.queueMu.Unlock()
	player1, player2 := m.queue[0], m.queue[1]
	m.queue = m.queue[2:]
	defer player1.Conn.Close()
	defer player2.Conn.Close()
	return player1.UserID, player2.UserID
}

func (m *MatchService) SearchPlayerConn() {
	var countUser int
	select {
	case _ = <-m.channel:
		if countUser == 2 {
			m.CreatePair()
			countUser = 0
		} else {
			countUser += 1
		}
	default:
		fmt.Println("Ничего нет в канале — работаем дальше")
	}
}
