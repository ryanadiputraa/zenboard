package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ryanadiputraa/zenboard/internal/domain"
)

type taskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) domain.TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (r *taskRepository) FetchTasks(ctx context.Context, boardID string) (tasks []domain.TaskDAO, err error) {
	err = r.db.Select(&tasks, `
		SELECT ts.id, ts.order, ts.name, ts.board_id,
		t.id AS task_id, t.order AS task_order, t.name AS task_name, t.tag, t.assignee, t.created_at, t.updated_at
		FROM task_status AS ts
		LEFT JOIN tasks AS t ON t.status_id = ts.id
		WHERE ts.board_id = $1
		ORDER BY ts.order ASC, t.order ASC
		`, boardID)
	if err != nil {
		return
	}

	return
}
