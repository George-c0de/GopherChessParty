package services

import (
	"GopherChessParty/ent"
	"GopherChessParty/ent/chess"
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/errors"
	"GopherChessParty/internal/interfaces"
	"encoding/json"
	"fmt"
	chesslib "github.com/corentings/chess/v2"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
type GameService struct {
	log        interfaces.ILogger
	repository interfaces.IGameRepo
	games      map[uuid.UUID]*Game
}

func (m *GameService) CreateGame(playerID1, playerID2 uuid.UUID) (*ent.Chess, error) {
	game, err := m.repository.Create(playerID1, playerID2)
	if err != nil {
		m.log.Error(err)
		return nil, err
	}
	m.startGame(game.ID, playerID1, playerID2)
	return game, nil
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

func (m *GameService) GetGame(GameID uuid.UUID, player *dto.PlayerConn) (*Game, error) {

	if m.games[GameID] == nil {
		status := m.repository.GetStatus(GameID)
		if status != chess.StatusInProgress {
			m.log.Error(errors.ErrGameEnd)
			return nil, errors.ErrGameEnd
		}
		m.GetStatusGame(GameID)
	}
	game := m.games[GameID]
	if player.UserID == game.WhitePlayer.UserID && game.WhitePlayer.Conn == nil {
		game.WhitePlayer = player
	} else if player.UserID == game.BlackPlayer.UserID && game.BlackPlayer.Conn == nil {
		game.BlackPlayer = player
	}
	return game, nil
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
	var err error
	game, err := m.GetGame(GameID, player)
	if err != nil {
		return err
	}
	if move == "" {
		return errors.ErrMoveEmpty
	}
	var currentMotionUser *dto.PlayerConn
	var oponentMotionUser *dto.PlayerConn
	var oponentMotion int
	if game.currentMotion == WhiteMotion {
		currentMotionUser = game.WhitePlayer
		oponentMotionUser = game.BlackPlayer
		oponentMotion = BlackMotion
	} else {
		currentMotionUser = game.BlackPlayer
		oponentMotionUser = game.WhitePlayer
		oponentMotion = WhiteMotion
	}

	if currentMotionUser.UserID != player.UserID {
		return errors.ErrCurrentUserMotion
	}
	err = game.Match.PushNotationMove(move, chesslib.UCINotation{}, &chesslib.PushMoveOptions{ForceMainline: true})
	if err != nil {
		m.log.Error(err)
		return err
	}
	fmt.Printf(game.Match.CurrentPosition().Board().Draw())
	fmt.Println(game.Match.Outcome())
	if game.Match.Outcome() != chesslib.NoOutcome {
		err := m.UpdateStatus(GameID, game.Match.Outcome())
		if err != nil {
			return err
		}
		return errors.ErrGameEnd
	}
	if oponentMotionUser.Conn != nil {
		err = m.sendMove(oponentMotionUser, move)
		if err != nil {
			return err
		}
	}

	game.currentMotion = oponentMotion
	return nil
}

// sendMove отправляет сообщение игроку через WebSocket
func (m *GameService) sendMove(player *dto.PlayerConn, move string) error {
	answer := map[string]interface{}{
		"move": move,
	}
	response, err := json.Marshal(answer)
	if err != nil {
		m.log.Error(err)
		return err
	}
	err = player.Conn.WriteMessage(websocket.TextMessage, response)
	if err != nil {
		m.log.Error(err)
		return err
	}
	return nil
}

func (m *GameService) SetPlayer(GameID uuid.UUID, player *dto.PlayerConn) {
	m.GetGame(GameID, player)
}
