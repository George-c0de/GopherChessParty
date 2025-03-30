package services

import (
	"log/slog"

	"GopherChessParty/internal/config"
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/models"
	"GopherChessParty/internal/storage/interfaces"
	"github.com/corentings/chess/v2"
	"github.com/google/uuid"
)

type IService interface {
	IUserService
	IGameService
	IAuthService
	CreateUser(data *dto.CreateUser) (*models.User, error)
	ValidPassword(data dto.AuthenticateUser) bool
}
type Service struct {
	IUserService
	IGameService
	IAuthService
	logger *slog.Logger
}

func New(log *slog.Logger, repository interfaces.IRepository, cfg config.Auth) *Service {
	userService := UserService{log: log, repository: repository}
	gameService := GameService{log: log, nowGames: map[uuid.UUID]*chess.Game{}}
	authService := AuthService{log: log, jwtSecret: cfg.JwtSecret, exp: cfg.ExpTime}
	return &Service{
		IUserService: &userService,
		IGameService: &gameService,
		IAuthService: &authService,
		logger:       log,
	}
}

func (s *Service) CreateUser(data *dto.CreateUser) (*models.User, error) {
	return s.IUserService.CreateUser(data)
}

func (s *Service) ValidPassword(data dto.AuthenticateUser) bool {
	HashedPassword, err := s.GetUserPassword(data.Email)
	if err != nil {
		s.logger.Error("Not Found User", err, data.Email)
	}
	return s.IsValidPassword(HashedPassword, data.Password)
}
