package services

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/models"
	"GopherChessParty/internal/storage/interfaces"
	"log/slog"
)

type IUserService interface {
	GetUsers() ([]*models.User, error)
	CreateUser(data *dto.CreateUser) (*models.User, error)
}

type UserService struct {
	repository interfaces.IRepository
	log        *slog.Logger
}

func (m *UserService) CreateUser(data *dto.CreateUser) (*models.User, error) {
	return m.repository.CreateUser(data)
}
func (m *UserService) GetUsers() ([]*models.User, error) {
	return m.repository.GetUsers()
}
