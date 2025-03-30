package services

import (
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/logger"
	"GopherChessParty/internal/models"
	"GopherChessParty/internal/storage/interfaces"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	GetUsers() ([]*models.User, error)
	CreateUser(data *dto.CreateUser) (*models.User, error)
	GetUserPassword(Email string) (string, error)
}

type UserService struct {
	repository interfaces.IRepository
	log        *logger.Logger
}

func (m *UserService) CreateUser(data *dto.CreateUser) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		m.log.Error(err)
		return nil, err
	}

	data.Password = string(hashedPassword)

	return m.repository.CreateUser(data)
}

func (m *UserService) GetUsers() ([]*models.User, error) {
	return m.repository.GetUsers()
}

func (m *UserService) GetUserPassword(email string) (string, error) {
	return m.repository.GetUserPassword(email)
}
