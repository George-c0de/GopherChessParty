package interfaces

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/models"
)

type IUserRepo interface {
	CreateUser(user *dto.CreateUser) (*models.User, error)
	GetUsers() ([]*models.User, error)
	GetUserPassword(email string) (*models.AuthUser, error)
}
