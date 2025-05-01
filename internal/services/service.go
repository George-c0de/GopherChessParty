package services

import (
	"GopherChessParty/ent"
	"GopherChessParty/internal/dto"
	custErr "GopherChessParty/internal/errors"
	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/models"
	"GopherChessParty/internal/utils"
	"errors"

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
	service := &Service{
		IUserService:  userService,
		IGameService:  gameService,
		IAuthService:  authService,
		IMatchService: matchService,
		logger:        logger,
	}
	go service.SearchPlayerConn()
	return service
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

func (s *Service) RegisterUser(data *dto.CreateUser) (*models.User, error) {
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

// CreateMatch создает пару игроков из очереди
func (s *Service) CreateMatch(playerID1, playerID2 uuid.UUID) (*ent.Chess, error) {
	game, err := s.CreateGame(playerID1, playerID2)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return game, nil
}

// SearchPlayerConn ищет пары игроков в очереди
func (s *Service) SearchPlayerConn() {
	for range s.IMatchService.GetExistsChannel() {
		if s.IMatchService.CheckPair() {
			player1, player2, err := s.IMatchService.ReturnPlayerID()
			if err != nil {
				s.logger.Error(err)
				continue
			}
			game, err := s.CreateMatch(player1.UserID, player2.UserID)
			if err != nil {
				s.logger.Error(err)
				continue
			}
			err = s.IMatchService.SendGemID(player1, game.ID)
			if err != nil {
				s.logger.Error(err)
				continue
			}
			err = s.IMatchService.SendGemID(player2, game.ID)
			if err != nil {
				s.logger.Error(err)
				continue
			}
			player1.Conn.Close()
			player2.Conn.Close()
		}
	}
}

func (s *Service) GetGameByID(gameID uuid.UUID) (*dto.Match, error) {
	return s.IGameService.GetGameByID(gameID)
}

func (s *Service) MoveGameStr(gameID uuid.UUID, move string, player *dto.PlayerConn) (int, bool) {
	err := s.IGameService.MoveGame(gameID, move, player)
	if err != nil {
		s.logger.Error(err)
		if errors.Is(err, custErr.ErrGameEnd) {
			return utils.GameFinished, false
		}
		return utils.GameInProgress, false
	}
	return utils.GameInProgress, true
}

func (s *Service) SetConnGame(GameID uuid.UUID, player *dto.PlayerConn) {
	s.IGameService.SetPlayer(GameID, player)
}
