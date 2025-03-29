package models

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type BaseUser struct {
	Name string
}

func (u BaseUser) String() string {
	return fmt.Sprintf("User(Name=%q)", u.Name)
}

type User struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u User) String() string {
	return fmt.Sprintf("User(Name=%q)", u.Name)
}
