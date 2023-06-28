package repository

import (
	"context"

	"github.com/ryanadiputraa/zenboard/domain"
	db "github.com/ryanadiputraa/zenboard/pkg/db/sqlc"
)

type userRepository struct {
	db *db.Queries
}

func NewUserRepository(db *db.Queries) domain.IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Save(ctx context.Context, user domain.User) (createdUser domain.User, err error) {
	return
}

func (r *userRepository) List(ctx context.Context, ids []string) (users []domain.User, err error) {
	return
}

func (r *userRepository) FindByID(ctx context.Context, userID string) (user domain.User, err error) {
	return
}
