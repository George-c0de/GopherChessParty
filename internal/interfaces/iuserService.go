package interfaces

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/models"
)

type IUserService interface {
	GetUsers() ([]*models.User, error)
	SaveUser(data *dto.CreateUser, hashedPassword string) (*models.User, error)
	GetUserPassword(Email string) (*models.AuthUser, error)
}
