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
		SELECT tasks.id, tasks.order, tasks.name, tasks.board_id,
		item.id AS item_id, item.order AS item_order, item.name AS item_name, item.tag, item.assignee, item.created_at, item.updated_at
		FROM tasks LEFT JOIN task_items AS item ON item.status_id = tasks.id
		WHERE tasks.board_id = $1
		ORDER BY tasks.order ASC, item.order ASC
		`, boardID)
	if err != nil {
		return
	}

	return
}
