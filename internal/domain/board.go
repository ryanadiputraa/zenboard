package domain

import (
	"context"
	"time"
)

type BoardRepository interface {
	Init(ctx context.Context, board InitBoardDTO, task1, task2, task3 InitTaskDTO) (Board, error)
	FetchByOwnerID(ctx context.Context, ownerID string) ([]Board, error)
}

type BoardService interface {
	InitBoard(ctx context.Context, userID string) error
	GetUserBoards(ctx context.Context, userID string) ([]Board, error)
}

type Board struct {
	ID          string    `json:"id" db:"id"`
	ProjectName string    `json:"project_name" db:"project_name"`
	Picture     string    `json:"picture" db:"picture"`
	OwnerID     string    `json:"owner_id" db:"owner_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
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
