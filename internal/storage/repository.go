package storage

import (
	"fmt"
	"github.com/George-c0de/GopherChessParty/internal/dto"
	"github.com/George-c0de/GopherChessParty/internal/models"
	"github.com/George-c0de/GopherChessParty/internal/storage/postgres"
)

type IRepository interface {
}

type Repository struct {
	*postgres.Postgres
}

func NewRepository(postgres *postgres.Postgres) *Repository {
	return &Repository{
		Postgres: postgres,
	}
}
func (r *Repository) GetUsers() []*models.User {
	var users []*models.User
	err := r.Db.Select(&users, "SELECT id, name, email, created_at, updated_at FROM users;")
	if err != nil {
		r.Log.Error("Ошибка при получении пользователей", err)
		return users
	}
	return users
}

func (r *Repository) CreateUser(user *dto.CreateUser) (*models.User, error) {
	query := `
		INSERT INTO users (name, email, password)
		VALUES (:name, :email, :password)
		RETURNING id, name, email, created_at, updated_at
	`
	rows, err := r.Db.NamedQuery(query, user)
	if err != nil {
		r.Log.Error("Ошибка при создании пользователя", err)
		return nil, err
	}
	defer rows.Close()

	var createdUser models.User
	if rows.Next() {
		if err := rows.StructScan(&createdUser); err != nil {
			r.Log.Error("Ошибка при сканировании созданного пользователя", err)
			return nil, err
		}
		return &createdUser, nil
	}

	return nil, fmt.Errorf("пользователь не создан")
}
