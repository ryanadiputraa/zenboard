package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ryanadiputraa/zenboard/internal/domain"
	"github.com/stretchr/testify/suite"
)

var (
	userID    = uuid.NewString()
	initBoard = domain.InitBoardDTO{
		ID:          uuid.NewString(),
		ProjectName: "Untitled",
		OwnerID:     userID,
		CreatedAt:   time.Now(),
	}
	task1 = domain.InitTaskDTO{
		ID:      uuid.NewString(),
		Order:   1,
		Name:    "Backlog",
		BoardID: initBoard.ID,
	}
	task2 = domain.InitTaskDTO{
		ID:      uuid.NewString(),
		Order:   2,
		Name:    "Do",
		BoardID: initBoard.ID,
	}
	task3 = domain.InitTaskDTO{
		ID:      uuid.NewString(),
		Order:   3,
		Name:    "Done",
		BoardID: initBoard.ID,
	}
)

type BoardRepositoryTestSuite struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository domain.BoardRepository
}

func TestBoardRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &BoardRepositoryTestSuite{})
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &BoardRepositoryTestSuite{})
}

func (ts *BoardRepositoryTestSuite) SetupTest() {
	conn, mock, err := sqlmock.New()
	if err != nil {
		ts.FailNow("fail to create mock db: %s", err)
	}
	mockDB := sqlx.NewDb(conn, "sqlmock")
	repository := NewBoardRepository(mockDB)

	ts.mock = mock
	ts.repository = repository
}

func (ts *BoardRepositoryTestSuite) TestInit() {
	cases := []struct {
		name              string
		expected          domain.Board
		err               error
		mockRepoBehaviour func(mock sqlmock.Sqlmock)
	}{
		{
			name: "should init user board with 3 default task status, if user didn't have board",
			expected: domain.Board{
				ID:          initBoard.ID,
				ProjectName: initBoard.ProjectName,
				Picture:     "",
				OwnerID:     initBoard.OwnerID,
				CreatedAt:   initBoard.CreatedAt,
			},
			err: nil,
			mockRepoBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id FROM boards WHERE owner_id").WillReturnError(sql.ErrNoRows)
				mock.ExpectQuery("^INSERT INTO boards").
					WithArgs(initBoard.ID, initBoard.ProjectName, "", initBoard.OwnerID, initBoard.CreatedAt).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "project_name", "picture", "owner_id", "created_at"}).AddRow(
							initBoard.ID, initBoard.ProjectName, "", initBoard.OwnerID, initBoard.CreatedAt,
						),
					)
				mock.ExpectExec(`INSERT INTO tasks`).WillReturnResult(
					sqlmock.NewResult(1, 3),
				)
				mock.ExpectCommit()
			},
		},
		{
			name:     "should rollback tx if err when fetch user board",
			expected: domain.Board{},
			err:      sql.ErrTxDone,
			mockRepoBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id FROM boards WHERE owner_id").WillReturnError(sql.ErrTxDone)
				mock.ExpectRollback()
			},
		},
		{
			name:     "should rollback tx if user already has a board",
			expected: domain.Board{},
			err:      nil,
			mockRepoBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id FROM boards WHERE owner_id").WillReturnRows(
					sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewString()),
				)
				mock.ExpectRollback()
			},
		},
		{
			name:     "should rollback tx if fail to insert into board table",
			expected: domain.Board{},
			err:      sql.ErrTxDone,
			mockRepoBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id FROM boards WHERE owner_id").WillReturnError(sql.ErrNoRows)
				mock.ExpectQuery("^INSERT INTO boards").
					WithArgs(initBoard.ID, initBoard.ProjectName, "", initBoard.OwnerID, initBoard.CreatedAt).
					WillReturnError(sql.ErrTxDone)
				mock.ExpectRollback()
			},
		},
		{
			name: "should rollback tx if fail to insert into task status table",
			expected: domain.Board{
				ID:          initBoard.ID,
				ProjectName: initBoard.ProjectName,
				Picture:     "",
				OwnerID:     initBoard.OwnerID,
				CreatedAt:   initBoard.CreatedAt,
			},
			err: sql.ErrTxDone,
			mockRepoBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT id FROM boards WHERE owner_id").WillReturnError(sql.ErrNoRows)
				mock.ExpectQuery("^INSERT INTO boards").
					WithArgs(initBoard.ID, initBoard.ProjectName, "", initBoard.OwnerID, initBoard.CreatedAt).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "project_name", "picture", "owner_id", "created_at"}).AddRow(
							initBoard.ID, initBoard.ProjectName, "", initBoard.OwnerID, initBoard.CreatedAt,
						),
					)
				mock.ExpectExec(`INSERT INTO tasks`).WillReturnError(sql.ErrTxDone)
				mock.ExpectRollback()
			},
		},
	}

	for _, c := range cases {
		ts.Run(c.name, func() {
			c.mockRepoBehaviour(ts.mock)
			created, err := ts.repository.Init(context.TODO(), initBoard, task1, task2, task3)
			ts.Equal(c.err, err)
			ts.Equal(c.expected, created)
		})
	}
}

func (ts *BoardRepositoryTestSuite) TestFetchByOwnerID() {
	ownerID := uuid.NewString()
	boards := []domain.Board{
		{
			ID:          uuid.NewString(),
			ProjectName: gofakeit.AppName(),
			Picture:     gofakeit.ImageURL(30, 30),
			OwnerID:     ownerID,
			CreatedAt:   gofakeit.Date(),
		},
		{
			ID:          uuid.NewString(),
			ProjectName: gofakeit.AppName(),
			Picture:     gofakeit.ImageURL(30, 30),
			OwnerID:     ownerID,
			CreatedAt:   gofakeit.Date(),
		},
		{
			ID:          uuid.NewString(),
			ProjectName: gofakeit.AppName(),
			Picture:     gofakeit.ImageURL(30, 30),
			OwnerID:     ownerID,
			CreatedAt:   gofakeit.Date(),
		},
	}

	cases := []struct {
		name              string
		expected          []domain.Board
		err               error
		mockRepoBehaviour func(mock sqlmock.Sqlmock)
	}{
		{
			name:     "should return list of user boards",
			expected: boards,
			err:      nil,
			mockRepoBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM boards").WillReturnRows(
					sqlmock.NewRows([]string{"id", "project_name", "picture", "owner_id", "created_at"}).
						AddRow(boards[0].ID, boards[0].ProjectName, boards[0].Picture, boards[0].OwnerID, boards[0].CreatedAt).
						AddRow(boards[1].ID, boards[1].ProjectName, boards[1].Picture, boards[1].OwnerID, boards[1].CreatedAt).
						AddRow(boards[2].ID, boards[2].ProjectName, boards[2].Picture, boards[2].OwnerID, boards[2].CreatedAt),
				)
			},
		},
		{
			name:     "should fail to return list of user boards",
			expected: nil,
			err:      sql.ErrNoRows,
			mockRepoBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM boards").WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, c := range cases {
		ts.Run(c.name, func() {
			c.mockRepoBehaviour(ts.mock)
			list, err := ts.repository.FetchByOwnerID(context.TODO(), ownerID)
			ts.Equal(c.err, err)
			ts.Equal(c.expected, list)
		})
	}
}
