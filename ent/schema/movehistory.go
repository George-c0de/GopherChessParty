package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// GameHistory holds the schema definition for the GameHistory entity.
type GameHistory struct {
	ent.Schema
}

// Fields of the GameHistory.
func (GameHistory) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.Time("created_at").Default(time.Now),
		field.Int("num"),
		field.String("move"),
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("game_id", uuid.UUID{}),
	}
}

// Edges of the GameHistory.
func (GameHistory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("moves").
			Field("user_id").
			Unique().
			Required(),
		edge.From("game", Chess.Type).
			Ref("moves").
			Field("game_id").
			Unique().
			Required(),
	}
}
