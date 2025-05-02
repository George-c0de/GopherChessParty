package storage

import (
	"GopherChessParty/ent"
	"GopherChessParty/ent/chess"
	"GopherChessParty/ent/user"
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"context"
	"github.com/google/uuid"
)

type GameRepository struct {
	log interfaces.ILogger
	*Repository
}

func NewGameRepository(log interfaces.ILogger, client *Repository) *GameRepository {
	return &GameRepository{
		log:        log,
		Repository: client,
	}
}

func (g *GameRepository) GetGames(userID uuid.UUID) ([]*dto.GameHistory, error) {
	ctx := context.Background()
	games, err := g.client.Chess.
		Query().
		Select(
			chess.FieldID,
			chess.FieldCreatedAt,
			chess.FieldResult,
			chess.FieldStatus,
			chess.BlackUserColumn,
			chess.WhiteUserColumn,
		).WithWhiteUser(func(uq *ent.UserQuery) {
		uq.Select(user.FieldID, user.FieldName)
	}).
		WithBlackUser(func(uq *ent.UserQuery) {
			uq.Select(user.FieldID, user.FieldName)
		}).
		Where(chess.Or(
			// партii, где userID — чёрный
			chess.HasBlackUserWith(user.IDEQ(userID)),
			// или партii, где userID — белый
			chess.HasWhiteUserWith(user.IDEQ(userID)),
		)).All(ctx)
	if err != nil {
		g.log.Error(err)
		return nil, err
	}
	gamesWithUsers := make([]*dto.GameHistory, 0, len(games))
	for _, game := range games {
		gamesWithUsers = append(gamesWithUsers, &dto.GameHistory{
			Id:        game.ID,
			CreatedAt: game.CreatedAt,
			UpdatedAt: game.UpdatedAt,
			Status:    game.Status,
			Result:    game.Result,
			BlackPlayer: &dto.Player{
				Id:   game.Edges.BlackUser.ID,
				Name: game.Edges.BlackUser.Name,
			},
			WhitePlayer: &dto.Player{
				Id:   game.Edges.WhiteUser.ID,
				Name: game.Edges.WhiteUser.Name,
			},
		})
	}
	return gamesWithUsers, nil
}

func (g *GameRepository) Create(playerID1, playerID2 uuid.UUID) (*ent.Chess, error) {
	ctx := context.Background()
	game, err := g.client.Chess.Create().
		SetWhiteUserID(playerID1).
		SetBlackUserID(playerID2).
		Save(ctx)
	if err != nil {
		g.log.Error(err)
		return nil, err
	}
	return game, nil
}

func (g *GameRepository) GetGameById(gameId uuid.UUID) (*dto.Match, error) {
	ctx := context.Background()
	game, err := g.client.Chess.Query().
		WithBlackUser().
		WithWhiteUser().
		Where(chess.ID(gameId)).
		Only(ctx)
	if err != nil {
		g.log.Error(err)
		return nil, err
	}
	return &dto.Match{
		ID:        game.ID,
		CreatedAt: game.CreatedAt,
		Status:    game.Status,
		Result:    game.Result,
		BlackUser: &dto.GetUser{
			ID:    game.Edges.BlackUser.ID,
			Name:  game.Edges.BlackUser.Name,
			Email: game.Edges.BlackUser.Email,
		},
		WhiteUser: &dto.GetUser{
			ID:    game.Edges.WhiteUser.ID,
			Name:  game.Edges.WhiteUser.Name,
			Email: game.Edges.WhiteUser.Email,
		},
	}, nil
}

func (g *GameRepository) GetStatus(GameID uuid.UUID) chess.Status {
	var status chess.Status
	ctx := context.Background()
	game, err := g.client.Chess.Query().Select(chess.FieldStatus).Where(chess.ID(GameID)).Only(ctx)
	if err != nil {
		g.log.Error(err)
		return status
	}
	return game.Status
}

func (g *GameRepository) UpdateGameResult(
	GameId uuid.UUID,
	status chess.Status,
	result chess.Result,
) error {
	ctx := context.Background()
	_, err := g.client.Chess.UpdateOneID(GameId).SetStatus(status).SetResult(result).Save(ctx)
	if err != nil {
		g.log.Error(err)
		return err
	}
	return nil
}

func (g *GameRepository) SaveMove(GameID uuid.UUID, move string, UserID uuid.UUID) (error, *ent.GameHistory) {
	ctx := context.Background()
	save, err := g.client.GameHistory.Create().SetGameID(GameID).SetMove(move).SetUserID(UserID).SetNum(2).Save(ctx)
	if err != nil {
		g.log.Error(err)
		return err, nil
	}
	return nil, save
}
