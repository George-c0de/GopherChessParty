package services

import (
	"GopherChessParty/ent"
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Service struct {
	interfaces.IUserService
	interfaces.IGameService
	interfaces.IAuthService
	interfaces.IMatchService
	logger interfaces.ILogger
}

func NewService(
	userService interfaces.IUserService,
	gameService interfaces.IGameService,
	authService interfaces.IAuthService,
	matchService interfaces.IMatchService,
	logger interfaces.ILogger,
) *Service {
	return &Service{
		IUserService:  userService,
		IGameService:  gameService,
		IAuthService:  authService,
		IMatchService: matchService,
		logger:        logger,
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

func (s *Service) GetGamesByUserID(userId uuid.UUID) ([]*ent.Chess, error) {
	return s.IGameService.GetGames(userId)
}
