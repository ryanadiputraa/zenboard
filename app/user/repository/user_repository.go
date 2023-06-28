package repository

import (
	"context"
	"database/sql"

	"github.com/ryanadiputraa/zenboard/domain"
	db "github.com/ryanadiputraa/zenboard/pkg/db/sqlc"
)

type userRepository struct {
	db *db.Queries
}

func NewUserRepository(db *db.Queries) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Save(ctx context.Context, user domain.User) (createdUser domain.User, err error) {
	arg := db.CreateUserParams{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.LastName,
		Picture:    sql.NullString{String: user.Picture, Valid: true},
		Locale:     user.Locale,
		BoardLimit: user.BoardLimit,
		CreatedAt:  user.CreatedAt,
	}

	created, err := r.db.CreateUser(ctx, arg)
	if err != nil {
		return
	}

	createdUser = domain.User{
		ID:         created.ID,
		FirstName:  created.FirstName,
		LastName:   created.LastName,
		Email:      created.Email,
		Picture:    created.Picture.String,
		Locale:     created.Locale,
		BoardLimit: created.BoardLimit,
		CreatedAt:  created.CreatedAt,
	}
	return
}

func (r *userRepository) List(ctx context.Context, ids []string) (users []domain.User, err error) {
	list, err := r.db.ListUsers(ctx, ids)
	if err != nil {
		return
	}

	for _, v := range list {
		user := domain.User{
			ID:         v.ID,
			FirstName:  v.FirstName,
			LastName:   v.LastName,
			Email:      v.Email,
			Picture:    v.Picture.String,
			Locale:     v.Locale,
			BoardLimit: v.BoardLimit,
			CreatedAt:  v.CreatedAt,
		}
		users = append(users, user)
	}
	return
}

func (r *userRepository) FindByID(ctx context.Context, userID string) (user domain.User, err error) {
	data, err := r.db.GetUser(ctx, userID)
	if err != nil {
		return
	}

	user = domain.User{
		ID:         data.ID,
		FirstName:  data.FirstName,
		LastName:   data.LastName,
		Email:      data.Email,
		Picture:    data.Picture.String,
		Locale:     data.Locale,
		BoardLimit: data.BoardLimit,
		CreatedAt:  data.CreatedAt,
	}
	return
}
