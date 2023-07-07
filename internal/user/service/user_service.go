package service

import (
	"context"
	"time"

	"github.com/ryanadiputraa/zenboard/internal/domain"
	log "github.com/sirupsen/logrus"
)

type userService struct {
	repository domain.UserRepository
}

func NewUserService(repository domain.UserRepository) domain.UserService {
	return &userService{
		repository: repository,
	}
}

func (s *userService) CreateOrUpdateUserIfExists(ctx context.Context, user domain.User) (createdUser domain.User, err error) {
	user.BoardLimit = domain.DEFAULT_BOARD_LIMIT
	user.CreatedAt = time.Now()

	createdUser, err = s.repository.Save(ctx, user)
	if err != nil {
		log.Error("fail to create or update user: ", err)
	}
	return
}

func (s *userService) FindUserByID(ctx context.Context, userID string) (user domain.User, err error) {
	user, err = s.repository.FindByID(ctx, userID)
	if err != nil {
		log.Info("fail to find user: ", err)
	}
	return
}
