package domain

import (
	"context"
	"database/sql"
)

type TaskRepository interface {
	ListByBoardID(ctx context.Context, boardID string) ([]TaskDAO, error)
	Create(ctx context.Context, task Task) (Task, error)
	DeleteByID(ctx context.Context, taskID string) (Task, error)
}

type TaskService interface {
	ListBoardTasks(ctx context.Context, boardID string) ([]TaskDTO, error)
	AddBoardTask(ctx context.Context, boardID, taskName string) (TaskDTO, error)
	DeleteTask(ctx context.Context, taskID string) error
}

type TaskItem struct {
	ID          string `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
	Order       int    `json:"order" db:"order"`
	Tag         string `json:"tag" db:"tag"`
	Assignee    string `json:"assignee" db:"assignee"`
	BoardID     string `json:"-" db:"board_id"`
	StatusID    string `json:"-" db:"status_id"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
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
	BoardID  string     `json:"-"`
	TaskItem []TaskItem `json:"tasks"`
}

type TaskDAO struct {
	ID              string         `db:"id"`
	Order           int            `db:"order"`
	Name            string         `db:"name"`
	BoardID         string         `db:"board_id"`
	ItemID          sql.NullString `db:"item_id"`
	ItemDescription sql.NullString `db:"item_description"`
	ItemOrder       sql.NullInt16  `db:"item_order"`
	Tag             sql.NullString `db:"tag"`
	Assignee        sql.NullString `db:"assignee"`
	StatusID        sql.NullString `db:"status_id"`
	CreatedAt       sql.NullString `db:"created_at"`
	UpdatedAt       sql.NullString `db:"updated_at"`
}

func GenerateTaskList(daoList []TaskDAO) (tasks []TaskDTO) {
	idx := -1
	taskMap := make(map[string]bool)

	for _, l := range daoList {
		if _, exists := taskMap[l.ID]; !exists {
			idx++
			var ts TaskDTO
			ts.ID = l.ID
			ts.Order = l.Order
			ts.Name = l.Name
			ts.BoardID = l.BoardID

			ts.TaskItem = []TaskItem{}
			tasks = append(tasks, ts)
			taskMap[ts.ID] = true
		}

		if !l.ItemID.Valid {
			continue
		}

		var t TaskItem
		t.ID = l.ItemID.String
		t.Description = l.ItemDescription.String
		t.Order = int(l.ItemOrder.Int16)
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
