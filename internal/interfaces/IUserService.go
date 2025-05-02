package interfaces

import (
	"GopherChessParty/internal/dto"
	"github.com/google/uuid"
)

type IUserService interface {
	GetUsers() ([]*dto.User, error)
	SaveUser(data *dto.CreateUser, hashedPassword string) (*dto.User, error)
	GetUserPassword(Email string) (*dto.AuthUser, error)
	GetUserByID(UserID uuid.UUID) (*dto.User, error)
}
