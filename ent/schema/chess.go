package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Chess holds the schema definition for the Chess entity.
type Chess struct {
	ent.Schema
}

const (
	Open = iota
	Games
	Closed
)

// Fields of the Chess.
func (Chess) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.Uint8("status").Default(Open),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now),
	}
}

// Edges of the Chess.
func (Chess) Edges() []ent.Edge {
	return []ent.Edge{}
}
