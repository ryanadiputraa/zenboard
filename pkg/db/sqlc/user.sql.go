// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  id, first_name, last_name, email, picture,
  locale, board_limit, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id) DO UPDATE SET 
  first_name = excluded.first_name,
  last_name = excluded.last_name,
  picture = excluded.picture,
  locale = excluded.locale
RETURNING id, first_name, last_name, email, picture, locale, board_limit, created_at
`

type CreateUserParams struct {
	ID         string         `json:"id"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	Email      string         `json:"email"`
	Picture    sql.NullString `json:"picture"`
	Locale     string         `json:"locale"`
	BoardLimit int32          `json:"board_limit"`
	CreatedAt  time.Time      `json:"created_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Picture,
		arg.Locale,
		arg.BoardLimit,
		arg.CreatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Picture,
		&i.Locale,
		&i.BoardLimit,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, first_name, last_name, email, picture, locale, board_limit, created_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Picture,
		&i.Locale,
		&i.BoardLimit,
		&i.CreatedAt,
	)
	return i, err
}
