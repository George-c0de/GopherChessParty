package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("email").MaxLen(100),
		field.String("name").MaxLen(255),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now),
		field.String("password").MaxLen(255),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
