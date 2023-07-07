
-- name: CreateBoard :one
INSERT INTO boards (
  id, project_name, picture, owner_id, created_at
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetBoardByOwnerID :one
SELECT * FROM boards
WHERE owner_id = $1 LIMIT 1;