package interfaces

import (
	"GopherChessParty/internal/dto"
	"github.com/google/uuid"
)

type IUserService interface {
	Users() ([]*dto.User, error)
	SaveUser(data *dto.CreateUser, hashedPassword string) (*dto.User, error)
	UserPassword(Email string) (*dto.AuthUser, error)
	UserByID(UserID uuid.UUID) (*dto.User, error)
}
