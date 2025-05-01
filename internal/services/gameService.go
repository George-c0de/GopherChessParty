package services

import (
	"GopherChessParty/ent"
	"GopherChessParty/ent/chess"
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/errors"
	"GopherChessParty/internal/interfaces"
	chesslib "github.com/corentings/chess/v2"
	"github.com/google/uuid"
)

// GameService логика для работы с игрой
type GameService struct {
	log        interfaces.ILogger
	repository interfaces.IGameRepo
	games      map[uuid.UUID]*dto.Game
}

func NewGameService(
	log interfaces.ILogger,
	repository interfaces.IGameRepo,
) interfaces.IGameService {
	return &GameService{
		log:        log,
		repository: repository,
		games:      make(map[uuid.UUID]*dto.Game),
	}
}

func (m *GameService) CreateGame(playerID1, playerID2 uuid.UUID) (*ent.Chess, error) {
	game, err := m.repository.Create(playerID1, playerID2)
	if err != nil {
		return nil, err
	}
	m.startGame(game.ID, playerID1, playerID2)
	return game, nil
}

func (m *GameService) GetGames(userID uuid.UUID) ([]*ent.Chess, error) {
	return m.repository.GetGames(userID)
}

func (m *GameService) GetGameByID(gameID uuid.UUID) (*dto.Match, error) {
	return m.repository.GetGameById(gameID)
}

func (m *GameService) startGame(GameID uuid.UUID, whiteUserID, blackUserID uuid.UUID) {
	m.games[GameID] = &dto.Game{
		Match:         chesslib.NewGame(),
		GameID:        GameID,
		CurrentMotion: dto.WhiteMotion,
		WhitePlayer:   &dto.PlayerConn{UserID: whiteUserID},
		BlackPlayer:   &dto.PlayerConn{UserID: blackUserID},
		HistoryMove:   make([]string, 0),
		Status:        chess.StatusInProgress,
		Result:        chess.Result00,
	}
}

func (m *GameService) GetStatusGame(GameID uuid.UUID) chess.Status {
	return m.repository.GetStatus(GameID)
}

func (m *GameService) GetGame(GameID uuid.UUID) (*dto.Game, error) {
	if m.games[GameID] == nil {
		status := m.repository.GetStatus(GameID)
		if status != chess.StatusInProgress {
			m.log.Error(errors.ErrGameEnd)
			return nil, errors.ErrGameEnd
		}
		m.log.Error(errors.ErrGameEnd)
		return nil, errors.ErrGameDeleteFailed
	}
	return m.games[GameID], nil
}

func (m *GameService) UpdateStatus(
	gameID uuid.UUID,
	parse chesslib.Outcome,
) error {
	status := chess.StatusFinished
	var result chess.Result
	switch parse {
	case chesslib.BlackWon:
		result = chess.Result01
	case chesslib.WhiteWon:
		result = chess.Result10
	case chesslib.Draw:
		result = chess.Result11
	default:
		status = chess.StatusAborted
		result = chess.Result00
	}
	game, err := m.GetGame(gameID)
	if err != nil {
		m.log.Error(err)
		return err
	}
	game.Status = status
	game.Result = result
	return m.repository.UpdateGameResult(gameID, status, result)
}

func (m *GameService) MoveValid(GameID uuid.UUID, move string) error {
	if move == "" {
		return errors.ErrInvalidMove
	}
	game, err := m.GetGame(GameID)
	if err != nil {
		return err
	}
	wanted, err := chesslib.UCINotation{}.Decode(game.Match.Position(), move)
	if err != nil {
		m.log.Error(err)
		return err
	}
	for _, move := range game.Match.ValidMoves() {
		if wanted.S1() == move.S1() && wanted.S2() == move.S2() {
			return nil
		}
	}
	return errors.ErrInvalidMove
}

func (m *GameService) MoveGame(GameID uuid.UUID, move string, player *dto.PlayerConn) error {
	game, err := m.GetGame(GameID)
	if err != nil {
		return err
	}

	currentMotionUser := game.GetCurrentUser()
	if currentMotionUser.UserID != player.UserID {
		m.log.Error(errors.ErrCurrentUserMotion)
		return errors.ErrCurrentUserMotion
	}

	err = game.Match.PushNotationMove(
		move,
		chesslib.UCINotation{},
		&chesslib.PushMoveOptions{},
	)
	if err != nil {
		m.log.Error(err)
		return err
	}
	game.SetMove(move)
	if game.Match.Outcome() != chesslib.NoOutcome {
		err := m.UpdateStatus(GameID, game.Match.Outcome())
		if err != nil {
			m.log.Error(err)
			return err
		}
		return errors.ErrGameEnd
	}
	return nil
}

func (m *GameService) SetPlayer(GameID uuid.UUID, player *dto.PlayerConn) error {
	game := m.games[GameID]
	if game == nil {
		m.log.Error(errors.ErrGameNotFound)
		return errors.ErrGameNotFound
	}
	if player.UserID == game.WhitePlayer.UserID {
		game.WhitePlayer.Conn = player.Conn
	} else if player.UserID == game.BlackPlayer.UserID {
		game.BlackPlayer.Conn = player.Conn
	}
	return nil
}

func (m *GameService) IsConnectPlayers(GameID uuid.UUID) bool {
	game := m.games[GameID]
	if game == nil {
		return false
	}
	return game.BlackPlayer != nil && game.BlackPlayer.Conn != nil && game.WhitePlayer != nil &&
		game.WhitePlayer.Conn != nil
}

func (m *GameService) GetOpponent(gameID uuid.UUID) *dto.PlayerConn {
	game := m.games[gameID]
	return game.GetOpponentUser()
}

func (m *GameService) GetStatusMemory(gameID uuid.UUID) (chess.Status, error) {
	game := m.games[gameID]
	if game == nil {
		return chess.StatusInProgress, errors.ErrGameNotFound
	}
	return game.Status, nil
}

func (m *GameService) GetHistoryMove(gameID uuid.UUID) []string {
	game := m.games[gameID]
	if game == nil {
		return []string{}
	}
	return game.HistoryMove
}

func (m *GameService) GetGameMemory(GameID uuid.UUID) *dto.Game {
	return m.games[GameID]
}
