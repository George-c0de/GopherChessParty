package services

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"github.com/google/uuid"
)

type UserService struct {
	repository interfaces.IUserRepo
	log        interfaces.ILogger
}

func NewUserService(
	logger interfaces.ILogger,
	repository interfaces.IUserRepo,
) interfaces.IUserService {
	return &UserService{repository: repository, log: logger}
}

func (m *UserService) SaveUser(data *dto.CreateUser, hashedPassword string) (*dto.User, error) {
	data.Password = hashedPassword
	return m.repository.CreateUser(data)
}

func (m *UserService) Users() ([]*dto.User, error) {
	return m.repository.Users()
}

func (m *UserService) UserPassword(email string) (*dto.AuthUser, error) {
	return m.repository.UserPassword(email)
}

func (m *UserService) UserByID(userID uuid.UUID) (*dto.User, error) {
	return m.repository.UserByID(userID)
}
