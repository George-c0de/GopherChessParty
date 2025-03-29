package services

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/models"
	"github.com/corentings/chess/v2"
)

type IService interface {
	CreateGame() *chess.Game
	GetUsers() ([]*models.User, error)
	CreateUser(data *dto.CreateUser) (*models.User, error)
}
