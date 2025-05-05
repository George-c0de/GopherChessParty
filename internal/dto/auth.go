package dto

import "github.com/google/uuid"

type AuthUser struct {
	UserID         uuid.UUID
	HashedPassword string
}
