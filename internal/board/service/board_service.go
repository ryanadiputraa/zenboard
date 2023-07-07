package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ryanadiputraa/zenboard/internal/domain"
	log "github.com/sirupsen/logrus"
)

type boardService struct {
	repository domain.BoardRepository
}

func NewBoardService(repository domain.BoardRepository) domain.BoardService {
	return &boardService{
		repository: repository,
	}
}

func (s *boardService) InitBoard(ctx context.Context, userID string) (err error) {
	board := domain.InitBoardDTO{
		ID:          uuid.NewString(),
		ProjectName: "Untitled",
		OwnerID:     userID,
		CreatedAt:   time.Now(),
	}

	task1 := domain.InitTaskDTO{
		ID:      uuid.NewString(),
		Order:   1,
		Name:    "Backlog",
		BoardID: board.ID,
	}
	task2 := domain.InitTaskDTO{
		ID:      uuid.NewString(),
		Order:   2,
		Name:    "Do",
		BoardID: board.ID,
	}
	task3 := domain.InitTaskDTO{
		ID:      uuid.NewString(),
		Order:   3,
		Name:    "Done",
		BoardID: board.ID,
	}

	created, err := s.repository.Init(ctx, board, task1, task2, task3)
	if err != nil {
		log.Error("fail to init user board: ", err)
		return
	}
	if created.ID != "" {
		log.WithFields(log.Fields{
			"id":           created.ID,
			"project_name": created.ProjectName,
			"owner_id":     created.OwnerID,
			"created_at":   created.CreatedAt,
		}).Info("init new user board")
	}

	return
}
