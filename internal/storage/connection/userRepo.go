package connection

import (
	"context"

	"GopherChessParty/ent/user"
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"GopherChessParty/internal/models"
)

type UserRepository struct {
	log interfaces.ILogger
	*Repository
}

func NewUserRepository(log interfaces.ILogger, repo *Repository) *UserRepository {
	return &UserRepository{log, repo}
}

func (r *UserRepository) CreateUser(user *dto.CreateUser) (*models.User, error) {
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

func (r *UserRepository) GetUsers() ([]*models.User, error) {
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

func (r *UserRepository) GetUserPassword(email string) (*models.AuthUser, error) {
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
