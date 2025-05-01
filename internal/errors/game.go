package errors

import "errors"

var (
	ErrGameEnd           = errors.New("game is over")
	ErrMoveEmpty         = errors.New("move is empty")
	ErrCurrentUserMotion = errors.New("not your motion")
	ErrGameNotFound      = errors.New("game not found")
	ErrGameDeleteFailed  = errors.New("game delete failed")
	ErrInvalidMove       = errors.New("invalid move")
	ErrPlayersNotConn    = errors.New("players not connected")
)
