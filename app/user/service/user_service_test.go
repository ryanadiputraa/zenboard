package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/ryanadiputraa/zenboard/domain"
	"github.com/ryanadiputraa/zenboard/pkg/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, &UserServiceTestSuite{})
}

func createRandomUser() domain.User {
	return domain.User{
		ID:            uuid.NewString(),
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		Email:         gofakeit.Email(),
		Picture:       gofakeit.ImageURL(120, 120),
		Locale:        gofakeit.Country(),
		BoardLimit:    domain.DEFAULT_BOARD_LIMIT,
		CreatedAt:     time.Now(),
		VerifiedEmail: true,
	}
}

func (ts *UserServiceTestSuite) TestCreateOrUpdateUserIfExists() {
	randomUser := createRandomUser()

	cases := []struct {
		name              string
		expected          domain.User
		err               error
		mockRepoBehaviour func(mockRepo *mocks.UserRepository)
	}{
		{
			name:     "should return created or updated user",
			expected: randomUser,
			err:      nil,
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("Save", mock.Anything, mock.Anything).Return(randomUser, nil)
			},
		},
		{
			name:     "should fail to create or update and return empty user and sql error",
			expected: domain.User{},
			err:      sql.ErrNoRows,
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("Save", mock.Anything, mock.Anything).Return(domain.User{}, sql.ErrNoRows)
			},
		},
	}

	for _, c := range cases {
		ts.Run(c.name, func() {
			mockRepo := mocks.UserRepository{}
			service := NewUserService(&mockRepo)
			c.mockRepoBehaviour(&mockRepo)

			res, err := service.CreateOrUpdateUserIfExists(context.TODO(), randomUser)
			ts.Equal(c.err, err)
			ts.Equal(c.expected, res)
		})
	}
}

func (ts *UserServiceTestSuite) TestFindUserByID() {
	randomUsers := createRandomUser()

	cases := []struct {
		name              string
		expected          domain.User
		err               error
		mockRepoBehaviour func(mockRepo *mocks.UserRepository)
	}{
		{
			name:     "should return a user with specified id",
			expected: randomUsers,
			err:      nil,
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByID", mock.Anything, mock.Anything).Return(randomUsers, nil)
			},
		},
		{
			name:     "should return empty user and sql error",
			expected: domain.User{},
			err:      sql.ErrNoRows,
			mockRepoBehaviour: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByID", mock.Anything, mock.Anything).Return(domain.User{}, sql.ErrNoRows)
			},
		},
	}

	for _, c := range cases {
		ts.Run(c.name, func() {
			mockRepo := mocks.UserRepository{}
			service := NewUserService(&mockRepo)
			c.mockRepoBehaviour(&mockRepo)

			res, err := service.FindUserByID(context.TODO(), randomUsers.ID)
			ts.Equal(c.err, err)
			ts.Equal(c.expected, res)
		})
	}
}
