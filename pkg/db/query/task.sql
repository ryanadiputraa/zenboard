-- name: InitTaskStatus :many
INSERT INTO task_status (
  id, "order", name, board_id
) VALUES
  ($2, $3, $4, $1),
  ($5, $6, $7, $1),
  ($8, $9, $10, $1)
RETURNING *;