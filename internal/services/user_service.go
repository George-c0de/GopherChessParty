package services

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/models"
)

type UserService struct {
	repository interfaces.IRepository
	log        interfaces.ILogger
}

func (m *UserService) SaveUser(
	data *dto.CreateUser,
	hashedPassword string,
) (*models.User, error) {
	data.Password = hashedPassword
	return m.repository.CreateUser(data)
}

func (m *UserService) GetUsers() ([]*models.User, error) {
	return m.repository.GetUsers()
}

func (m *UserService) GetUserPassword(email string) (*models.AuthUser, error) {
	return m.repository.GetUserPassword(email)
}
