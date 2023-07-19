package service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
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

func (s *taskService) ListBoardTasks(ctx context.Context, boardID string) (tasks []domain.TaskDTO, err error) {
	list, err := s.repository.ListByBoardID(ctx, boardID)
	if err != nil && err != sql.ErrNoRows {
		log.Error("fail to fetch task list: ", err)
	}

	tasks = domain.GenerateTaskList(list)
	return
}

func (s *taskService) AddBoardTask(ctx context.Context, boardID, taskName string) (created domain.Task, err error) {
	task := domain.Task{
		ID:      uuid.NewString(),
		Name:    taskName,
		BoardID: boardID,
	}

	created, err = s.repository.Create(ctx, task)
	if err != nil {
		log.Error("fail to create task: ", err)
		return
	}
	log.WithFields(log.Fields{
		"id":        created.ID,
		"task_name": created.Name,
		"order":     created.Order,
		"board_id":  created.BoardID,
	}).Info("created task")

	return
}
