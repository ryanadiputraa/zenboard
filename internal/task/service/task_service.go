package service

import (
	"context"
	"database/sql"

	"github.com/ryanadiputraa/zenboard/internal/domain"
	log "github.com/sirupsen/logrus"
)

type taskService struct {
	repository domain.TaskRepository
}

func NewTaskService(repository domain.TaskRepository) domain.TaskService {
	return &taskService{
		repository: repository,
	}
}

func (s *taskService) ListBoardTasks(ctx context.Context, boardID string) (tasks []domain.TaskStatus, err error) {
	list, err := s.repository.FetchTasks(ctx, boardID)
	if err != nil && err != sql.ErrNoRows {
		log.Error("fail to fetch task list: ", err)
	}

	tasks = domain.GenerateTaskList(list)

	return
}
