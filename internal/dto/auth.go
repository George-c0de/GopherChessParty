package dto

import "github.com/google/uuid"

type AuthUser struct {
	UserId         uuid.UUID
	HashedPassword string
}
