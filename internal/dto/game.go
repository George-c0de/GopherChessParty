package dto

import (
	"GopherChessParty/ent/chess"
	"github.com/google/uuid"
	"time"
)

type Player struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
type GameHistory struct {
	Id          uuid.UUID    `json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Status      chess.Status `json:"status"`
	Result      chess.Result `json:"result"`
	BlackPlayer *Player      `json:"black_player"`
	WhitePlayer *Player      `json:"white_player"`
}
