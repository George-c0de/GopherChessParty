package dto

import (
	"GopherChessParty/ent/chess"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"time"
)

type PlayerConn struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
}

type Match struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Result    chess.Result
	Status    chess.Status
	WhiteUser *GetUser
	BlackUser *GetUser
}
