package services

import (
	"GopherChessParty/ent"
	"GopherChessParty/ent/chess"
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/errors"
	"GopherChessParty/internal/interfaces"
	chesslib "github.com/corentings/chess/v2"
	"github.com/google/uuid"
	"time"
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

func (m *GameService) GamesByUserID(userID uuid.UUID) ([]*dto.GameHistory, error) {
	return m.repository.Games(userID)
}

func (m *GameService) GameByID(gameID uuid.UUID) (*dto.Match, error) {
	return m.repository.GameById(gameID)
}

func (m *GameService) startGame(GameID uuid.UUID, whiteUserID, blackUserID uuid.UUID) {
	m.games[GameID] = &dto.Game{
		Match:         chesslib.NewGame(),
		CreatedAt:     time.Now(),
		ID:            GameID,
		CurrentMotion: dto.WhiteMotion,
		WhitePlayer:   &dto.PlayerConn{UserID: whiteUserID},
		BlackPlayer:   &dto.PlayerConn{UserID: blackUserID},
		HistoryMove:   make([]string, 0),
		Status:        chess.StatusInProgress,
		Result:        chess.Result00,
	}
}

func (m *GameService) StatusGame(GameID uuid.UUID) chess.Status {
	return m.repository.Status(GameID)
}

func (m *GameService) Game(GameID uuid.UUID) (*dto.Game, error) {
	return m.GameDB(GameID)
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
	game, err := m.Game(gameID)
	if err != nil {
		m.log.Error(err)
		return err
	}
	game.Status = status
	game.Result = result
	return m.repository.UpdateGame(gameID, status, result)
}

func (m *GameService) MoveValid(GameID uuid.UUID, move string) error {
	if move == "" {
		return errors.ErrInvalidMove
	}
	game, err := m.Game(GameID)
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
	game, err := m.Game(GameID)
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
	_, err = m.repository.SaveMove(GameID, move, player.UserID, game.NumMove)
	if err != nil {
		return err
	}
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
	game, err := m.GameDB(GameID)
	if err != nil {
		return err
	}
	switch player.UserID {
	case game.WhitePlayer.UserID:
		game.WhitePlayer.Conn = player.Conn
	case game.BlackPlayer.UserID:
		game.BlackPlayer.Conn = player.Conn
	default:
		return errors.ErrPlayerNotFound
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

func (m *GameService) Opponent(gameID uuid.UUID) *dto.PlayerConn {
	game := m.games[gameID]
	return game.GetOpponentUser()
}

func (m *GameService) GameMemory(GameID uuid.UUID) *dto.Game {
	return m.games[GameID]
}
func (m *GameService) GameDB(gameID uuid.UUID) (*dto.Game, error) {
	game := m.GameMemory(gameID)
	if game != nil {
		return game, nil
	}
	gameDB, err := m.GameByID(gameID)
	if err != nil {
		return nil, err
	}
	match := chesslib.NewGame()
	currentMotion := dto.WhiteMotion
	historyMove := make([]string, 0, len(gameDB.HistoryMove))
	NumMoves := 0
	for _, move := range gameDB.HistoryMove {
		_ = match.PushNotationMove(
			move.Move,
			chesslib.UCINotation{},
			&chesslib.PushMoveOptions{},
		)
		NumMoves++
		if currentMotion == dto.WhiteMotion {
			currentMotion = dto.BlackMotion
		} else {
			currentMotion = dto.WhiteMotion
		}
		historyMove = append(historyMove, move.Move)
	}
	game = &dto.Game{
		ID:            gameID,
		CreatedAt:     gameDB.CreatedAt,
		Result:        gameDB.Result,
		Status:        gameDB.Status,
		Match:         match,
		WhitePlayer:   &dto.PlayerConn{UserID: gameDB.WhiteUser.ID},
		BlackPlayer:   &dto.PlayerConn{UserID: gameDB.BlackUser.ID},
		CurrentMotion: currentMotion,
		HistoryMove:   historyMove,
		NumMove:       NumMoves,
	}
	m.games[gameID] = game
	return game, nil
}
