package services

import (
	exc "errors"

	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/errors"
	"GopherChessParty/internal/interfaces"
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
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, false
	}
	return token, token.Valid
}

func (s *Service) CreateUser(data *dto.CreateUser) (*dto.User, error) {
	hashedPassword, err := s.GeneratePassword(data.Password)
	if err != nil {
		return nil, err
	}
	return s.SaveUser(data, hashedPassword)
}

func (s *Service) RegisterUser(data *dto.CreateUser) (*dto.User, error) {
	hashedPassword, err := s.GeneratePassword(data.Password)
	if err != nil {
		return nil, err
	}
	return s.SaveUser(data, hashedPassword)
}

func (s *Service) ValidPassword(data dto.AuthenticateUser) (*uuid.UUID, bool) {
	userAuth, err := s.UserPassword(data.Email)
	if err != nil {
		return nil, false
	}
	return &userAuth.UserID, s.IsValidPassword(userAuth.HashedPassword, data.Password)
}

// SearchPlayerConn ищет пары игроков в очереди
func (s *Service) SearchPlayerConn() {
	for range s.ExistsChannel() {
		if s.CheckPair() {
			player1, player2 := s.ReturnPlayers()
			game, err := s.CreateGame(player1.UserID, player2.UserID)
			if err != nil {
				_ = s.AddUser(player1)
				_ = s.AddUser(player2)
				continue
			}
			_ = s.SendGemID(player1, game.ID)
			_ = s.SendGemID(player2, game.ID)
			_ = player1.Conn.Close()
			_ = player2.Conn.Close()
		}
	}
}

func (s *Service) PlayerExit(player *dto.PlayerConn) error {
	s.ExitPlayerAdd(player.UserID)
	err := player.Conn.Close()
	if err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}

func (s *Service) MoveGameStr(gameID uuid.UUID, move string, player *dto.PlayerConn) bool {
	if !s.IsConnectPlayers(gameID) {
		s.logger.Error(errors.ErrPlayersNotConn)
		return false
	}

	opponentMotionUser := s.GetOpponent(gameID)

	err := s.MoveValid(gameID, move)
	if err != nil {
		return false
	}

	errMove := s.MoveGame(gameID, move, player)
	if errMove != nil && !exc.Is(errMove, errors.ErrGameEnd) {
		return false
	}

	response, errSend := s.GetGameInfoMemory(gameID, errMove == nil, move)
	if errSend != nil {
		return false
	}
	err = s.SendMessage(opponentMotionUser, response)
	return err == nil
}

func (s *Service) SetConnGame(GameID uuid.UUID, player *dto.PlayerConn) error {
	err := s.SetPlayer(GameID, player)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetGameInfoMemory(
	GameID uuid.UUID,
	ok bool,
	move string,
) (map[string]interface{}, error) {
	game := s.GameMemory(GameID)
	if game == nil {
		s.logger.Error(errors.ErrGameNotFound)
		return nil, errors.ErrGameNotFound
	}
	answer := map[string]interface{}{
		"ok":          ok,
		"status":      game.Status,
		"historyMove": game.HistoryMove,
		"currentMove": game.CurrentMotion,
		"result":      game.Result,
	}
	if !ok {
		answer["message"] = "Недопустимый ход"
	}
	if move != "" {
		answer["move"] = move
	}

	return answer, nil
}
