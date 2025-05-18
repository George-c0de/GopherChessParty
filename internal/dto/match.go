package dto

import (
	"time"

	"GopherChessParty/ent/chess"
	chess2 "github.com/corentings/chess/v2"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	WhiteMotion = iota
	BlackMotion
)

type PlayerConn struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
}

type Move struct {
	ID        uuid.UUID `json:"id"         db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Num       int       `json:"num"        db:"num"`
	Move      string    `json:"move"       db:"move"`
	UserID    uuid.UUID `json:"user_id"    db:"user_id"`
}

type Match struct {
	ID          uuid.UUID    `json:"id"           db:"id"`
	CreatedAt   time.Time    `json:"created_at"   db:"created_at"`
	Result      chess.Result `json:"result"       db:"result"`
	Status      chess.Status `json:"status"       db:"status"`
	WhiteUser   *GetUser     `json:"white_user"   db:"white_user"`
	BlackUser   *GetUser     `json:"black_user"   db:"black_user"`
	HistoryMove []*Move      `json:"history_move"`
}

type Game struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Result    chess.Result
	Status    chess.Status

	Match         *chess2.Game
	WhitePlayer   *PlayerConn
	BlackPlayer   *PlayerConn
	CurrentMotion int
	HistoryMove   []string
	NumMove       int
}

func (game *Game) GetOpponentUser() *PlayerConn {
	if game.CurrentMotion == WhiteMotion {
		return game.BlackPlayer
	} else {
		return game.WhitePlayer
	}
}

func (game *Game) GetCurrentUser() *PlayerConn {
	if game.CurrentMotion == WhiteMotion {
		return game.WhitePlayer
	} else {
		return game.BlackPlayer
	}
}

func (game *Game) SetMove(move string) {
	if game.CurrentMotion == WhiteMotion {
		game.CurrentMotion = BlackMotion
	} else {
		game.CurrentMotion = WhiteMotion
	}
	game.HistoryMove = append(game.HistoryMove, move)
	game.NumMove++
}
