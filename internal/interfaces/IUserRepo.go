package interfaces

import (
	"GopherChessParty/internal/dto"
	"github.com/google/uuid"
)

type IUserRepo interface {
	CreateUser(user *dto.CreateUser) (*dto.User, error)
	Users() ([]*dto.User, error)
	UserPassword(email string) (*dto.AuthUser, error)
	UserByID(UserID uuid.UUID) (*dto.User, error)
}
