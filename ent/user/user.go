// Code generated by ent, DO NOT EDIT.

package user

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldPassword holds the string denoting the password field in the database.
	FieldPassword = "password"
	// EdgeWhiteID holds the string denoting the white_id edge name in mutations.
	EdgeWhiteID = "white_id"
	// EdgeBlackID holds the string denoting the black_id edge name in mutations.
	EdgeBlackID = "black_id"
	// EdgeMoves holds the string denoting the moves edge name in mutations.
	EdgeMoves = "moves"
	// Table holds the table name of the user in the database.
	Table = "users"
	// WhiteIDTable is the table that holds the white_id relation/edge.
	WhiteIDTable = "chesses"
	// WhiteIDInverseTable is the table name for the Chess entity.
	// It exists in this package in order to avoid circular dependency with the "chess" package.
	WhiteIDInverseTable = "chesses"
	// WhiteIDColumn is the table column denoting the white_id relation/edge.
	WhiteIDColumn = "user_white_id"
	// BlackIDTable is the table that holds the black_id relation/edge.
	BlackIDTable = "chesses"
	// BlackIDInverseTable is the table name for the Chess entity.
	// It exists in this package in order to avoid circular dependency with the "chess" package.
	BlackIDInverseTable = "chesses"
	// BlackIDColumn is the table column denoting the black_id relation/edge.
	BlackIDColumn = "user_black_id"
	// MovesTable is the table that holds the moves relation/edge.
	MovesTable = "game_histories"
	// MovesInverseTable is the table name for the GameHistory entity.
	// It exists in this package in order to avoid circular dependency with the "gamehistory" package.
	MovesInverseTable = "game_histories"
	// MovesColumn is the table column denoting the moves relation/edge.
	MovesColumn = "user_id"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldEmail,
	FieldName,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldPassword,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// EmailValidator is a validator for the "email" field. It is called by the builders before save.
	EmailValidator func(string) error
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// PasswordValidator is a validator for the "password" field. It is called by the builders before save.
	PasswordValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the User queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByEmail orders the results by the email field.
func ByEmail(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEmail, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByPassword orders the results by the password field.
func ByPassword(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPassword, opts...).ToFunc()
}

// ByWhiteIDCount orders the results by white_id count.
func ByWhiteIDCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newWhiteIDStep(), opts...)
	}
}

// ByWhiteID orders the results by white_id terms.
func ByWhiteID(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newWhiteIDStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByBlackIDCount orders the results by black_id count.
func ByBlackIDCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newBlackIDStep(), opts...)
	}
}

// ByBlackID orders the results by black_id terms.
func ByBlackID(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newBlackIDStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByMovesCount orders the results by moves count.
func ByMovesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newMovesStep(), opts...)
	}
}

// ByMoves orders the results by moves terms.
func ByMoves(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMovesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newWhiteIDStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(WhiteIDInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, WhiteIDTable, WhiteIDColumn),
	)
}
func newBlackIDStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(BlackIDInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, BlackIDTable, BlackIDColumn),
	)
}
func newMovesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MovesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, MovesTable, MovesColumn),
	)
}
