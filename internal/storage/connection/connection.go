package connection

import (
	"context"

	"GopherChessParty/ent"
	"GopherChessParty/ent/user"
	"GopherChessParty/internal/config"
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/logger"
	"GopherChessParty/internal/models"
	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

type Repository struct {
	log    *logger.Logger
	client *ent.Client
}

// MustNewRepository Создание нового подключения.
func MustNewRepository(cfg config.Database, log *logger.Logger) *Repository {
	drv, err := sql.Open("postgres", cfg.DBUrl())
	if err != nil {
		panic(err)
	}

	db := drv.DB()
	db.SetMaxIdleConns(cfg.MaxIdleConns)   // Максимальное количество простаивающих соединений.
	db.SetMaxOpenConns(cfg.MaxOpenConns)   // Максимальное количество открытых соединений.
	db.SetConnMaxLifetime(cfg.MaxTimeLife) // Максимальное время жизни одного соединения.

	return &Repository{log: log, client: ent.NewClient(ent.Driver(drv))}
}

func (r *Repository) CreateUser(user *dto.CreateUser) (*models.User, error) {
	ctx := context.Background()

	newUser, err := r.client.User.Create().
		SetEmail(user.Email).
		SetPassword(user.Password).
		SetName(user.Name).
		Save(ctx)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	return &models.User{
		ID:        newUser.ID,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Name:      newUser.Name,
	}, nil
}

func (r *Repository) GetUsers() ([]*models.User, error) {
	ctx := context.Background()

	var users []*models.User
	err := r.client.User.
		Query().
		Select(user.FieldID, user.FieldName, user.FieldEmail, user.FieldCreatedAt, user.FieldUpdatedAt).
		Scan(ctx, &users)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserPassword(email string) (*models.AuthUser, error) {
	ctx := context.Background()

	authUser, err := r.client.User.
		Query().
		Where(user.Email(email)).
		Select(user.FieldPassword, user.FieldID).
		First(ctx)

	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return &models.AuthUser{UserId: authUser.ID, HashedPassword: authUser.Password}, nil
}
