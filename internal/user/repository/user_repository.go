package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ryanadiputraa/zenboard/internal/domain"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Save(ctx context.Context, arg domain.User) (user domain.User, err error) {
	err = r.db.QueryRowxContext(ctx, `INSERT INTO users (
			id, first_name, last_name, email, picture,
			locale, board_limit, created_at, verified_email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (email) DO UPDATE SET 
			first_name = excluded.first_name,
			last_name = excluded.last_name,
			picture = excluded.picture,
			locale = excluded.locale,
			board_limit = excluded.board_limit,
			verified_email = excluded.verified_email
		RETURNING *`,
		arg.ID, arg.FirstName, arg.LastName, arg.Email, arg.Picture,
		arg.Locale, arg.BoardLimit, arg.CreatedAt, arg.VerifiedEmail,
	).StructScan(&user)
	return
}

func (r *userRepository) FindByID(ctx context.Context, userID string) (user domain.User, err error) {
	err = r.db.Get(&user, "SELECT * FROM users WHERE id = $1 LIMIT 1", userID)
	return
}
