-- name: CreateUser :one
INSERT INTO users (
  id, first_name, last_name, email, picture,
  locale, board_limit, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id) DO UPDATE SET 
  first_name = excluded.first_name,
  last_name = excluded.last_name,
  picture = excluded.picture,
  locale = excluded.locale
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;