package repository

import (
	"context"
	"database/sql"

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

func (r *taskRepository) ListByBoardID(ctx context.Context, boardID string) (tasks []domain.TaskDAO, err error) {
	err = r.db.Select(&tasks, `
		SELECT tasks.id, tasks.order, tasks.name, tasks.board_id,
		item.id AS item_id, item.order AS item_order, item.description AS item_description, item.tag, item.assignee, item.created_at, item.updated_at
		FROM tasks LEFT JOIN task_items AS item ON item.status_id = tasks.id
		WHERE tasks.board_id = $1
		ORDER BY tasks.order ASC, item.order ASC
	`, boardID)
	return
}

func (r *taskRepository) Create(ctx context.Context, task domain.Task) (created domain.Task, err error) {
	err = r.db.QueryRowxContext(ctx, `
		INSERT INTO tasks (id, "order", name, board_id)
		VALUES ($1, (SELECT tasks."order" + 1 FROM tasks WHERE board_id = $3 ORDER BY tasks."order" DESC LIMIT 1), $2, $3)
		RETURNING *
	`, task.ID, task.Name, task.BoardID).StructScan(&created)
	return
}

func (r *taskRepository) DeleteByID(ctx context.Context, taskID string) (task domain.Task, err error) {
	err = r.db.QueryRowxContext(ctx, "DELETE FROM tasks WHERE id = $1 RETURNING *", taskID).StructScan(&task)
	return
}

func (r *taskRepository) UpdateOrder(ctx context.Context, task []domain.TaskReorderDTO) (err error) {
	tx := r.db.MustBeginTx(ctx, &sql.TxOptions{})
	for _, t := range task {
		_, err = r.db.ExecContext(ctx, `UPDATE tasks SET "order" = $2 WHERE id = $1`, t.ID, t.Order)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	err = tx.Commit()
	return
}
