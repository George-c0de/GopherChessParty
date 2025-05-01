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

const (
	WhiteMotion = iota
	BlackMotion
)

type Game struct {
	Match         *chesslib.Game
	WhitePlayer   *dto.PlayerConn
	BlackPlayer   *dto.PlayerConn
	currentMotion int
	GameID        uuid.UUID
}

func (game *Game) getOpponentUser() *dto.PlayerConn {
	if game.currentMotion == WhiteMotion {
		return game.BlackPlayer
	} else {
		return game.WhitePlayer
	}
}
func (game *Game) getCurrentUser() *dto.PlayerConn {
	if game.currentMotion == WhiteMotion {
		return game.WhitePlayer
	} else {
		return game.BlackPlayer
	}
}
func (game *Game) setMove() {
	if game.currentMotion == WhiteMotion {
		game.currentMotion = BlackMotion
	} else {
		game.currentMotion = WhiteMotion
	}
}

// GameService логика для работы с игрой
type GameService struct {
	log        interfaces.ILogger
	repository interfaces.IGameRepo
	games      map[uuid.UUID]*Game
}

func NewGameService(
	log interfaces.ILogger,
	repository interfaces.IGameRepo,
) interfaces.IGameService {
	return &GameService{
		log:        log,
		repository: repository,
		games:      make(map[uuid.UUID]*Game),
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
	m.games[GameID] = &Game{
		Match:         chesslib.NewGame(),
		GameID:        GameID,
		currentMotion: WhiteMotion,
		WhitePlayer:   &dto.PlayerConn{UserID: whiteUserID},
		BlackPlayer:   &dto.PlayerConn{UserID: blackUserID},
	}
}

func (m *GameService) GetStatusGame(GameID uuid.UUID) chess.Status {
	return m.repository.GetStatus(GameID)
}

func (m *GameService) GetGame(GameID uuid.UUID) (*Game, error) {
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

	return m.repository.UpdateGameResult(gameID, status, result)
}

func (m *GameService) MoveGame(GameID uuid.UUID, move string, player *dto.PlayerConn) error {
	game, err := m.GetGame(GameID)
	if err != nil {
		return err
	}
	if move == "" {
		return errors.ErrMoveEmpty
	}
	currentMotionUser := game.getCurrentUser()

	if currentMotionUser.UserID != player.UserID {
		return errors.ErrCurrentUserMotion
	}
	err = game.Match.PushNotationMove(
		move,
		chesslib.UCINotation{},
		&chesslib.PushMoveOptions{ForceMainline: true},
	)
	if err != nil {
		m.log.Error(err)
		return err
	}
	if game.Match.Outcome() != chesslib.NoOutcome {
		err := m.UpdateStatus(GameID, game.Match.Outcome())
		if err != nil {
			return err
		}
		return errors.ErrGameEnd
	}
	game.setMove()
	return nil
}

func (m *GameService) SetPlayer(GameID uuid.UUID, player *dto.PlayerConn) error {
	game := m.games[GameID]
	if game == nil {
		return errors.ErrGameNotFound
	}
	if player.UserID == game.WhitePlayer.UserID && game.WhitePlayer.Conn == nil {
		game.WhitePlayer.Conn = player.Conn
	} else if player.UserID == game.BlackPlayer.UserID && game.BlackPlayer.Conn == nil {
		game.BlackPlayer.Conn = player.Conn
	}
	return nil
}

func (m *GameService) IsConnectPlayers(GameID uuid.UUID) bool {
	game := m.games[GameID]
	if game == nil {
		return false
	}
	return game.BlackPlayer != nil && game.BlackPlayer.Conn != nil && game.WhitePlayer != nil && game.WhitePlayer.Conn != nil
}

func (m *GameService) GetOpponent(gameID uuid.UUID) *dto.PlayerConn {
	game := m.games[gameID]
	return game.getOpponentUser()
}
