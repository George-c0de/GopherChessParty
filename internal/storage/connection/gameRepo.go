package connection

import (
	"context"

	"GopherChessParty/ent"
	"GopherChessParty/ent/chess"
	"GopherChessParty/ent/user"
	"GopherChessParty/internal/interfaces"
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

func (g *GameRepository) GetGames(userID uuid.UUID) ([]*ent.Chess, error) {
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
		).
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
	return games, nil
}
func (g *GameRepository) Create(playerID1, playerID2 uuid.UUID) (*ent.Chess, error) {
	ctx := context.Background()
	game, err := g.client.Chess.Create().SetWhiteUserID(playerID1).SetBlackUserID(playerID2).Save(ctx)
	if err != nil {
		return nil, err
	}
	return game, nil
}
