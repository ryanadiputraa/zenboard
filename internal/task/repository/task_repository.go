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

func (r *taskRepository) FetchTasks(ctx context.Context, boardID string) (tasks []domain.TaskStatus, err error) {
	err = r.db.Select(&tasks, "SELECT * FROM task_status WHERE board_id = $1", boardID)
	if err != nil {
		return
	}

	for i, t := range tasks {
		var taskList []domain.Task

		err = r.db.Select(&taskList, "SELECT * FROM tasks WHERE status_id = $1", t.ID)
		if err != nil {
			return
		}

		if taskList == nil {
			tasks[i].Tasks = []domain.Task{}
		} else {
			tasks[i].Tasks = taskList
		}
	}

	return
}
