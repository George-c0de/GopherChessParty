package interfaces

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/models"
)

type IUserService interface {
	GetUsers() ([]*models.User, error)
	CreateUser(data *dto.CreateUser, hashedPassword string) (*models.User, error)
	GetUserPassword(Email string) (string, error)
}
