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
	ID        uuid.UUID
	CreatedAt time.Time
	Num       int
	Move      string
	UserID    uuid.UUID
}
type Match struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	Result      chess.Result
	Status      chess.Status
	WhiteUser   *GetUser
	BlackUser   *GetUser
	HistoryMove []*Move
}

type Game struct {
	Match         *chess2.Game
	WhitePlayer   *PlayerConn
	BlackPlayer   *PlayerConn
	CurrentMotion int
	GameID        uuid.UUID
	HistoryMove   []string
	Status        chess.Status
	Result        chess.Result
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
