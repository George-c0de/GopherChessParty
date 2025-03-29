package services

import (
	"GopherChessParty/internal/storage/interfaces"
	"log/slog"
)

type Service struct {
	IUserService
	IGameService
}

func New(log *slog.Logger, repository interfaces.IRepository) *Service {
	userService := UserService{log: log, repository: repository}
	gameService := GameService{log: log}
	return &Service{IUserService: &userService, IGameService: &gameService}
}
