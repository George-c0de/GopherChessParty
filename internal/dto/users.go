package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CreateUser struct {
	Name     string `binding:"required"       json:"name"`
	Email    string `binding:"required,email" json:"email"`
	Password string `binding:"required,min=8" json:"password"`
}

type AuthenticateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetUser struct {
	ID    uuid.UUID `json:"id"    db:"id"`
	Name  string    `json:"name"  db:"name"`
	Email string    `json:"email" db:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u User) String() string {
	return fmt.Sprintf("User(Name=%q)", u.Name)
}
