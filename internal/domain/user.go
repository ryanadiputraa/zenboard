package domain

import (
	"context"
	"time"
)

const (
	DEFAULT_BOARD_LIMIT = 3
)

type UserRepository interface {
	Save(ctx context.Context, user User) (User, error)
	FindByID(ctx context.Context, userID string) (User, error)
}

type UserService interface {
	CreateOrUpdateUserIfExists(ctx context.Context, user User) (User, error)
	FindUserByID(ctx context.Context, userID string) (User, error)
}

type User struct {
	ID            string    `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	Picture       string    `json:"picture"`
	Locale        string    `json:"locale"`
	BoardLimit    int32     `json:"board_limit"`
	CreatedAt     time.Time `json:"created_at"`
	VerifiedEmail bool      `json:"verified_email"`
}
