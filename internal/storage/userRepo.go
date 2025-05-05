package storage

import (
	"context"

	"GopherChessParty/ent/user"
	"GopherChessParty/internal/dto"
	"GopherChessParty/internal/interfaces"
	"github.com/google/uuid"
)

type UserRepository struct {
	log interfaces.ILogger
	*Repository
}

func NewUserRepository(log interfaces.ILogger, repo *Repository) *UserRepository {
	return &UserRepository{log, repo}
}

func (r *UserRepository) CreateUser(user *dto.CreateUser) (*dto.User, error) {
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
	return &dto.User{
		ID:        newUser.ID,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Name:      newUser.Name,
	}, nil
}

func (r *UserRepository) Users() ([]*dto.User, error) {
	ctx := context.Background()

	var users []*dto.User
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

func (r *UserRepository) UserByID(UserID uuid.UUID) (*dto.User, error) {
	ctx := context.Background()

	userDB, err := r.client.User.
		Query().
		Select(
			user.FieldID, user.FieldName, user.FieldEmail, user.FieldCreatedAt, user.FieldUpdatedAt,
		).Where(user.ID(UserID)).Only(ctx)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	return &dto.User{
		ID:        userDB.ID,
		Email:     userDB.Email,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
		Name:      userDB.Name,
	}, nil
}

func (r *UserRepository) UserPassword(email string) (*dto.AuthUser, error) {
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

	return &dto.AuthUser{UserID: authUser.ID, HashedPassword: authUser.Password}, nil
}
