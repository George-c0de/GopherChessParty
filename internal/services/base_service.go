package services

import (
	"GopherChessParty/internal/storage/interfaces"
	"log/slog"
)

type BaseService struct {
	log        *slog.Logger
	Repository interfaces.IRepository
}
