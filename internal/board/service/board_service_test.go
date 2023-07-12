package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/ryanadiputraa/zenboard/internal/domain"
	"github.com/ryanadiputraa/zenboard/pkg/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BoardSeriveTestSuite struct {
	suite.Suite
}

func TestBoardServiceTestSuite(t *testing.T) {
	suite.Run(t, &BoardSeriveTestSuite{})
}

func (ts *BoardSeriveTestSuite) TestInitBoard() {
	userID := uuid.NewString()

	cases := []struct {
		name              string
		err               error
		mockRepoBehaviour func(mockRepo *mocks.BoardRepository)
	}{
		{
			name: "should successfully init user board",
			err:  nil,
			mockRepoBehaviour: func(mockRepo *mocks.BoardRepository) {
				mockRepo.On("Init", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(domain.Board{
						ID:          uuid.NewString(),
						ProjectName: "Untitled",
						Picture:     "",
						OwnerID:     userID,
						CreatedAt:   time.Now(),
					}, nil)
			},
		},
		{
			name: "should fail to init user board",
			err:  sql.ErrTxDone,
			mockRepoBehaviour: func(mockRepo *mocks.BoardRepository) {
				mockRepo.On("Init", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(domain.Board{}, sql.ErrTxDone)
			},
		},
	}

	for _, c := range cases {
		ts.Run(c.name, func() {
			mockRepo := mocks.BoardRepository{}
			service := NewBoardService(&mockRepo)
			c.mockRepoBehaviour(&mockRepo)

			err := service.InitBoard(context.TODO(), userID)
			ts.Equal(c.err, err)
		})
	}
}

func (ts *BoardSeriveTestSuite) TestGetUserBoards() {
	userID := uuid.NewString()
	boards := []domain.Board{
		{
			ID:          uuid.NewString(),
			ProjectName: gofakeit.AppName(),
			Picture:     gofakeit.ImageURL(30, 30),
			OwnerID:     userID,
			CreatedAt:   gofakeit.Date(),
		},
		{
			ID:          uuid.NewString(),
			ProjectName: gofakeit.AppName(),
			Picture:     gofakeit.ImageURL(30, 30),
			OwnerID:     userID,
			CreatedAt:   gofakeit.Date(),
		},
		{
			ID:          uuid.NewString(),
			ProjectName: gofakeit.AppName(),
			Picture:     gofakeit.ImageURL(30, 30),
			OwnerID:     userID,
			CreatedAt:   gofakeit.Date(),
		},
	}

	cases := []struct {
		name              string
		expected          []domain.Board
		err               error
		mockRepoBehaviour func(mockRepo *mocks.BoardRepository)
	}{
		{
			name:     "should return list of board",
			expected: boards,
			err:      nil,
			mockRepoBehaviour: func(mockRepo *mocks.BoardRepository) {
				mockRepo.On("FetchByOwnerID", mock.Anything, mock.Anything).
					Return(boards, nil)
			},
		},
		{
			name:     "should fail to retrieve to list of user board",
			expected: []domain.Board{},
			err:      sql.ErrTxDone,
			mockRepoBehaviour: func(mockRepo *mocks.BoardRepository) {
				mockRepo.On("FetchByOwnerID", mock.Anything, mock.Anything).
					Return(nil, sql.ErrTxDone)
			},
		},
	}

	for _, c := range cases {
		ts.Run(c.name, func() {
			mockRepo := mocks.BoardRepository{}
			service := NewBoardService(&mockRepo)
			c.mockRepoBehaviour(&mockRepo)

			list, err := service.GetUserBoards(context.TODO(), userID)
			ts.Equal(c.err, err)
			ts.Equal(c.expected, list)
		})
	}
}
