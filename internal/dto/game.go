package dto

import (
	"time"

	"GopherChessParty/ent/chess"
	"github.com/google/uuid"
)

type Player struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
type GameHistory struct {
	ID          uuid.UUID    `json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Status      chess.Status `json:"status"`
	Result      chess.Result `json:"result"`
	BlackPlayer *Player      `json:"black_player"`
	WhitePlayer *Player      `json:"white_player"`
}
