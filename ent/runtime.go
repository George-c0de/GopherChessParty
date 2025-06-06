// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"GopherChessParty/ent/chess"
	"GopherChessParty/ent/gamehistory"
	"GopherChessParty/ent/schema"
	"GopherChessParty/ent/user"
	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	chessFields := schema.Chess{}.Fields()
	_ = chessFields
	// chessDescCreatedAt is the schema descriptor for created_at field.
	chessDescCreatedAt := chessFields[1].Descriptor()
	// chess.DefaultCreatedAt holds the default value on creation for the created_at field.
	chess.DefaultCreatedAt = chessDescCreatedAt.Default.(func() time.Time)
	// chessDescUpdatedAt is the schema descriptor for updated_at field.
	chessDescUpdatedAt := chessFields[2].Descriptor()
	// chess.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	chess.DefaultUpdatedAt = chessDescUpdatedAt.Default.(func() time.Time)
	// chessDescID is the schema descriptor for id field.
	chessDescID := chessFields[0].Descriptor()
	// chess.DefaultID holds the default value on creation for the id field.
	chess.DefaultID = chessDescID.Default.(func() uuid.UUID)
	gamehistoryFields := schema.GameHistory{}.Fields()
	_ = gamehistoryFields
	// gamehistoryDescCreatedAt is the schema descriptor for created_at field.
	gamehistoryDescCreatedAt := gamehistoryFields[1].Descriptor()
	// gamehistory.DefaultCreatedAt holds the default value on creation for the created_at field.
	gamehistory.DefaultCreatedAt = gamehistoryDescCreatedAt.Default.(func() time.Time)
	// gamehistoryDescID is the schema descriptor for id field.
	gamehistoryDescID := gamehistoryFields[0].Descriptor()
	// gamehistory.DefaultID holds the default value on creation for the id field.
	gamehistory.DefaultID = gamehistoryDescID.Default.(func() uuid.UUID)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescEmail is the schema descriptor for email field.
	userDescEmail := userFields[1].Descriptor()
	// user.EmailValidator is a validator for the "email" field. It is called by the builders before save.
	user.EmailValidator = userDescEmail.Validators[0].(func(string) error)
	// userDescName is the schema descriptor for name field.
	userDescName := userFields[2].Descriptor()
	// user.NameValidator is a validator for the "name" field. It is called by the builders before save.
	user.NameValidator = userDescName.Validators[0].(func(string) error)
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userFields[3].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescUpdatedAt is the schema descriptor for updated_at field.
	userDescUpdatedAt := userFields[4].Descriptor()
	// user.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	user.DefaultUpdatedAt = userDescUpdatedAt.Default.(func() time.Time)
	// userDescPassword is the schema descriptor for password field.
	userDescPassword := userFields[5].Descriptor()
	// user.PasswordValidator is a validator for the "password" field. It is called by the builders before save.
	user.PasswordValidator = userDescPassword.Validators[0].(func(string) error)
	// userDescID is the schema descriptor for id field.
	userDescID := userFields[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() uuid.UUID)
}
