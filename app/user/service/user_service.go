package service

import (
	"context"

	"github.com/ryanadiputraa/zenboard/domain"
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
	return
}

func (s *userService) ListUserWithinIds(ctx context.Context, ids []string) (users []domain.User, err error) {
	return
}

func (s *userService) FindUserByID(ctx context.Context, userID string) (user domain.User, err error) {
	return
}
