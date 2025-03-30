package postgres

import (
	"fmt"

	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	*Postgres
}

func NewRepository(postgres *Postgres) *Repository {
	return &Repository{
		Postgres: postgres,
	}
}

func (r *Repository) GetUsers() ([]*models.User, error) {
	var users []*models.User
	err := r.Db.Select(&users, "SELECT id, name, email, created_at, updated_at FROM users;")
	if err != nil {
		r.Log.Error("Ошибка при получении пользователей", err)
		return users, err
	}
	return users, nil
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
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			r.Log.Error(err.Error())
		}
	}(rows)

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
