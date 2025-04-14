package interfaces

import "github.com/corentings/chess/v2"

type IGameService interface {
	CreateGame() *chess.Game
	Move(game *chess.Game, move string) (bool, error)
}
