package services

import (
	"GopherChessParty/ent"
	"GopherChessParty/ent/chess"
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/errors"
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
		return nil, false
	}
	return token, token.Valid
}

func (s *Service) CreateUser(data *dto.CreateUser) (*models.User, error) {
	hashedPassword, err := s.IAuthService.GeneratePassword(data.Password)
	if err != nil {
		return nil, err
	}
	return s.IUserService.SaveUser(data, hashedPassword)
}

func (s *Service) RegisterUser(data *dto.CreateUser) (*models.User, error) {
	hashedPassword, err := s.IAuthService.GeneratePassword(data.Password)
	if err != nil {
		return nil, err
	}
	return s.IUserService.SaveUser(data, hashedPassword)
}

func (s *Service) ValidPassword(data dto.AuthenticateUser) (*uuid.UUID, bool) {
	userAuth, err := s.IUserService.GetUserPassword(data.Email)
	if err != nil {
		return nil, false
	}
	return &userAuth.UserId, s.IsValidPassword(userAuth.HashedPassword, data.Password)
}

func (s *Service) GetGamesByUserID(userId uuid.UUID) ([]*ent.Chess, error) {
	return s.IGameService.GetGames(userId)
}

// CreateMatch создает пару игроков из очереди
func (s *Service) CreateMatch(playerID1, playerID2 uuid.UUID) (*ent.Chess, error) {
	return s.CreateGame(playerID1, playerID2)
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
			// TODO(GEORGE): Сделать красивее
			player1.Conn.Close()
			player2.Conn.Close()
		}
	}
}

func (s *Service) GetGameByID(gameID uuid.UUID) (*dto.Match, error) {
	return s.IGameService.GetGameByID(gameID)
}

func (s *Service) MoveGameStr(gameID uuid.UUID, move string, player *dto.PlayerConn) error {
	if !s.IGameService.IsConnectPlayers(gameID) {
		s.logger.Error(errors.ErrPlayersNotConn)
		return errors.ErrPlayersNotConn
	}

	opponentMotionUser := s.IGameService.GetOpponent(gameID)

	err := s.ValidationMove(gameID, move)
	if err != nil {
		return err
	}

	errMove := s.IGameService.MoveGame(gameID, move, player)
	if errMove != nil {
		return errMove
	}

	errSend := s.IMatchService.SendMove(opponentMotionUser, move)
	if errSend != nil {
		return errSend
	}
	return nil
}

func (s *Service) SetConnGame(GameID uuid.UUID, player *dto.PlayerConn) {
	err := s.IGameService.SetPlayer(GameID, player)
	if err != nil {
		s.logger.Error(err)
		return
	}
}

func (s *Service) GetStatusMemory(GameID uuid.UUID) (chess.Status, error) {
	return s.IGameService.GetStatusMemory(GameID)
}

func (s *Service) GetHistoryMove(GameID uuid.UUID) []string {
	return s.IGameService.GetHistoryMove(GameID)
}

func (s *Service) GetGameInfoMemory(GameID uuid.UUID, ok bool) (map[string]interface{}, error) {
	game := s.IGameService.GetGameMemory(GameID)
	if game == nil {
		s.logger.Error(errors.ErrGameNotFound)
		return nil, errors.ErrGameNotFound
	}
	answer := map[string]interface{}{
		"ok":          ok,
		"status":      game.Status,
		"historyMove": game.HistoryMove,
		"currentMove": game.CurrentMotion,
	}
	if !ok {
		answer["message"] = "Недопустимый ход"
	}

	return answer, nil
}

func (s *Service) ValidationMove(GameID uuid.UUID, move string) error {
	return s.IGameService.MoveValid(GameID, move)
}
