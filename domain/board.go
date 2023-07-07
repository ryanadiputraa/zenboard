package domain

import (
	"context"
	"time"
)

type BoardRepository interface {
	Init(ctx context.Context, board InitBoardDTO, task1, task2, task3 InitTaskDTO) (Board, error)
}

type BoardService interface {
	InitBoard(ctx context.Context, userID string) error
}

type Board struct {
	ID          string    `json:"id"`
	ProjectName string    `json:"project_name"`
	Picture     string    `json:"picture"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type InitBoardDTO struct {
	ID          string
	ProjectName string
	OwnerID     string
	CreatedAt   time.Time
}

type InitTaskDTO struct {
	ID      string
	Order   int
	Name    string
	BoardID string
}
