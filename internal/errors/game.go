package errors

import "errors"

var (
	ErrGameEnd           = errors.New("game is over")
	ErrMoveEmpty         = errors.New("move is empty")
	ErrCurrentUserMotion = errors.New("not your motion")
)
