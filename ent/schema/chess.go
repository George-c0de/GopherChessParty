package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Chess holds the schema definition for the Chess entity.
type Chess struct {
	ent.Schema
}

const (
	waiting    = "waiting"
	inProgress = "in_progress"
	finished   = "finished"
	aborted    = "aborted"
)

const (
	winWhite   = "1-0"
	winBlack   = "0-1"
	draw       = "1-1"
	processing = "0-0"
)

// Fields of the Chess.
func (Chess) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now),
		field.Enum("status").Values(waiting, inProgress, finished, aborted).Default(waiting),
		field.Enum("result").Values(winWhite, winBlack, draw, processing).Default(processing),
	}
}

// Edges of the Chess.
func (Chess) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("white_user", User.Type).
			Ref("white_id").
			Unique().
			Required(),
		edge.From("black_user", User.Type).
			Ref("black_id").
			Unique().
			Required(),
		edge.To("moves", GameHistory.Type),
	}
}
