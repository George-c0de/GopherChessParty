package models

import (
	"fmt"
	"time"
)

type BaseUser struct {
	Name string
}

func (u BaseUser) String() string {
	return fmt.Sprintf("User(Name=%q)", u.Name)
}

type User struct {
	*BaseUser
	Id        int
	Email     string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Password  string
}
