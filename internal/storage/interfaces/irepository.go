package interfaces

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/models"
)

type IRepository interface {
	CreateUser(user *dto.CreateUser) (*models.User, error)
	GetUsers() ([]*models.User, error)
}
