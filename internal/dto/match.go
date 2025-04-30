package dto

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type PlayerConn struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
}
