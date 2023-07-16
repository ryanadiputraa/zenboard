package domain

import (
	"context"
	"database/sql"
)

type TaskRepository interface {
	FetchTasks(ctx context.Context, boardID string) ([]TaskDAO, error)
	// AddTask(ctx context.Context)
}

type TaskService interface {
	ListBoardTasks(ctx context.Context, boardID string) ([]TaskDTO, error)
}

type TaskItem struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Order     int    `json:"order" db:"order"`
	Tag       string `json:"tag" db:"tag"`
	Assignee  string `json:"assignee" db:"assignee"`
	BoardID   string `json:"-" db:"board_id"`
	StatusID  string `json:"-" db:"status_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type Task struct {
	ID      string `json:"id" db:"id"`
	Order   int    `json:"order" db:"order"`
	Name    string `json:"name" db:"name"`
	BoardID string `json:"board_id" db:"board_id"`
}

type TaskDTO struct {
	ID       string     `json:"id"`
	Order    int        `json:"order"`
	Name     string     `json:"name"`
	BoardID  string     `json:"board_id"`
	TaskItem []TaskItem `json:"tasks"`
}

type TaskDAO struct {
	ID        string         `db:"id"`
	Order     int            `db:"order"`
	Name      string         `db:"name"`
	BoardID   string         `db:"board_id"`
	TaskID    sql.NullString `db:"item_id"`
	TaskName  sql.NullString `db:"item_name"`
	TaskOrder sql.NullInt16  `db:"item_order"`
	Tag       sql.NullString `db:"tag"`
	Assignee  sql.NullString `db:"assignee"`
	StatusID  sql.NullString `db:"status_id"`
	CreatedAt sql.NullString `db:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at"`
}

func GenerateTaskList(daoList []TaskDAO) (tasks []TaskDTO) {
	idx := 0
	taskMap := make(map[string]bool)

	for _, l := range daoList {
		if _, exists := taskMap[l.ID]; !exists {
			var ts TaskDTO
			ts.ID = l.ID
			ts.Order = l.Order
			ts.Name = l.Name
			ts.BoardID = l.BoardID

			ts.TaskItem = []TaskItem{}
			tasks = append(tasks, ts)
			taskMap[ts.ID] = true
		}

		if !l.TaskID.Valid {
			continue
		}

		var t TaskItem
		t.ID = l.TaskID.String
		t.Name = l.TaskName.String
		t.Order = int(l.TaskOrder.Int16)
		t.Tag = l.Tag.String
		t.Assignee = l.Assignee.String
		t.StatusID = l.StatusID.String
		t.CreatedAt = l.CreatedAt.String
		t.UpdatedAt = l.UpdatedAt.String

		tasks[idx].TaskItem = append(tasks[idx].TaskItem, t)
	}

	if tasks == nil {
		tasks = []TaskDTO{}
	}

	return
}
