package interfaces

import (
	"GopherChessParty/internal/dto"
	"github.com/google/uuid"
)

type IUserRepo interface {
	CreateUser(user *dto.CreateUser) (*dto.User, error)
	GetUsers() ([]*dto.User, error)
	GetUserPassword(email string) (*dto.AuthUser, error)
	GetUserByID(UserID uuid.UUID) (*dto.User, error)
}
