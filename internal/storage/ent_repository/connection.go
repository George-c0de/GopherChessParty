package ent_repository

import (
	"GopherChessParty/ent"
	"GopherChessParty/ent/user"
	"GopherChessParty/internal/config"
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/models"
	"context"
	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

type Repository struct {
	client *ent.Client
}

func (r Repository) CreateUser(user *dto.CreateUser) (*models.User, error) {
	ctx := context.Background()
	newUser, err := r.client.User.Create().SetEmail(user.Email).SetPassword(user.Password).SetName(user.Name).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Id:        newUser.ID,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Name:      newUser.Name,
	}, nil
}

func (r Repository) GetUsers() ([]*models.User, error) {
	ctx := context.Background()
	var users []*models.User
	err := r.client.User.
		Query().
		Select(user.FieldID, user.FieldName, user.FieldEmail, user.FieldCreatedAt, user.FieldUpdatedAt).
		Scan(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func MustNewRepository(cfg config.Database) *Repository {
	drv, err := sql.Open("postgres", cfg.DBUrl())
	if err != nil {
		panic(err)
	}

	// Получаем underlying *sql.DB для настройки пула соединений.
	db := drv.DB()
	db.SetMaxIdleConns(cfg.MaxIdleConns)   // Максимальное количество простаивающих соединений.
	db.SetMaxOpenConns(cfg.MaxOpenConns)   // Максимальное количество открытых соединений.
	db.SetConnMaxLifetime(cfg.MaxTimeLife) // Максимальное время жизни одного соединения.

	return &Repository{ent.NewClient(ent.Driver(drv))}
}
