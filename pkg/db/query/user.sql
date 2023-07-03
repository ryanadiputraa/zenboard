-- name: CreateUser :one
INSERT INTO users (
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
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;