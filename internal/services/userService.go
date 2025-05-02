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

func (m *UserService) GetUsers() ([]*dto.User, error) {
	return m.repository.GetUsers()
}

func (m *UserService) GetUserPassword(email string) (*dto.AuthUser, error) {
	return m.repository.GetUserPassword(email)
}

func (m *UserService) GetUserByID(userID uuid.UUID) (*dto.User, error) {
	return m.repository.GetUserByID(userID)
}
