package domain

import "context"

type TaskRepository interface {
	FetchTasks(ctx context.Context, boardID string) ([]TaskStatus, error)
}

type TaskService interface {
	ListBoardTasks(ctx context.Context, boardID string) ([]TaskStatus, error)
}

type Task struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Order     int    `json:"order" db:"order"`
	Tag       string `json:"tag" db:"tag"`
	Assignee  string `json:"assignee" db:"assignee"`
	BoardID   string `json:"board_id" db:"board_id"`
	StatusID  string `json:"status_id" db:"status_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type TaskStatus struct {
	ID      string `json:"id" db:"id"`
	Order   int    `json:"order" db:"order"`
	Name    string `json:"name" db:"name"`
	BoardID string `json:"board_id" db:"board_id"`
	Tasks   []Task `json:"tasks"`
}
