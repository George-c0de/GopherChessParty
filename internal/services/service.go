package services

import (
	"GopherChessParty/internal/config"
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/logger"
	"GopherChessParty/internal/models"
	"github.com/corentings/chess/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Service struct {
	interfaces.IUserService
	interfaces.IGameService
	interfaces.IAuthService
	logger interfaces.ILogger
}

func NewService(log *logger.Logger, repository interfaces.IRepository, cfg config.Auth) *Service {
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

func (s *Service) IsValidateToken(tokenString string) (*jwt.Token, bool) {
	token, err := s.IAuthService.ValidateToken(tokenString)
	if err != nil {
		s.logger.Error(err)
		return nil, false
	}
	return token, token.Valid
}

func (s *Service) CreateUser(data *dto.CreateUser) (*models.User, error) {
	hashedPassword, err := s.IAuthService.GeneratePassword(data.Password)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return s.IUserService.SaveUser(data, hashedPassword)
}

func (s *Service) ValidPassword(data dto.AuthenticateUser) (*uuid.UUID, bool) {
	userAuth, err := s.IUserService.GetUserPassword(data.Email)
	if err != nil {
		s.logger.Error(err)
		return nil, false
	}
	return &userAuth.UserId, s.IsValidPassword(userAuth.HashedPassword, data.Password)
}

func (s *Service) GetGames(userId uuid.UUID) {

}
