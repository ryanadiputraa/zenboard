package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/ryanadiputraa/zenboard/domain"
	db "github.com/ryanadiputraa/zenboard/pkg/db/sqlc"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository domain.UserRepository
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &UserRepositoryTestSuite{})
}

func (ts *UserRepositoryTestSuite) SetupTest() {
	conn, mock, err := sqlmock.New()
	if err != nil {
		ts.FailNow("fail to create mock db: %s", err)
	}
	mockDB := db.New(conn)
	repository := NewUserRepository(mockDB)

	ts.mock = mock
	ts.repository = repository
}

func createRandomUser() domain.User {
	return domain.User{
		ID:         uuid.NewString(),
		FirstName:  gofakeit.FirstName(),
		LastName:   gofakeit.LastName(),
		Email:      gofakeit.Email(),
		Picture:    gofakeit.ImageURL(120, 120),
		Locale:     gofakeit.Country(),
		BoardLimit: domain.DEFAULT_BOARD_LIMIT,
		CreatedAt:  time.Now(),
	}
}

func (ts *UserRepositoryTestSuite) TestSave() {
	cases := []struct {
		name              string
		mockRepoBehaviour func(mock sqlmock.Sqlmock)
	}{
		{
			name: "should create a user",
			mockRepoBehaviour: func(mock sqlmock.Sqlmock) {
				user := createRandomUser()

				mock.ExpectQuery("INSERT INTO users").WillReturnRows(
					sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "picture", "locale", "board_limit", "created_at"}).
						AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.Picture, user.Locale, user.BoardLimit, user.CreatedAt),
				)

				created, err := ts.repository.Save(context.TODO(), user)
				ts.NoError(err)
				ts.NotEmpty(created)

				ts.Equal(user.ID, created.ID)
				ts.Equal(user.FirstName, created.FirstName)
				ts.Equal(user.LastName, created.LastName)
				ts.Equal(user.Email, created.Email)
				ts.Equal(user.Picture, created.Picture)
				ts.Equal(user.Locale, created.Locale)
				ts.Equal(user.BoardLimit, created.BoardLimit)
				ts.NotZero(created.CreatedAt)
			},
		},
		{
			name: "should fail to create user",
			mockRepoBehaviour: func(mock sqlmock.Sqlmock) {
				user := createRandomUser()

				mock.ExpectQuery("INSERT INTO users").WillReturnError(sql.ErrNoRows)

				created, err := ts.repository.Save(context.TODO(), user)
				ts.EqualError(sql.ErrNoRows, err.Error())
				ts.Empty(created)
			},
		},
	}

	for _, c := range cases {
		ts.Run(c.name, func() {
			c.mockRepoBehaviour(ts.mock)
		})
	}
}

func (ts *UserRepositoryTestSuite) TestFindByID() {
	cases := []struct {
		name              string
		mockRepoBehaviour func(mock sqlmock.Sqlmock)
	}{
		{
			name: "should return a user data with correct id",
			mockRepoBehaviour: func(mock sqlmock.Sqlmock) {
				user := createRandomUser()

				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(
					sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "picture", "locale", "board_limit", "created_at"}).
						AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.Picture, user.Locale, user.BoardLimit, user.CreatedAt),
				)

				data, err := ts.repository.FindByID(context.TODO(), user.ID)
				ts.NotEmpty(data)
				ts.NoError(err)

				ts.Equal(user.ID, data.ID)
				ts.Equal(user.FirstName, data.FirstName)
				ts.Equal(user.LastName, data.LastName)
				ts.Equal(user.Email, data.Email)
				ts.Equal(user.Picture, data.Picture)
				ts.Equal(user.Locale, data.Locale)
				ts.Equal(user.BoardLimit, data.BoardLimit)
				ts.NotZero(data.CreatedAt)
			},
		},
		{
			name: "should fail to find user that didn't exists",
			mockRepoBehaviour: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnError(sql.ErrNoRows)

				user, err := ts.repository.FindByID(context.TODO(), uuid.NewString())
				ts.Empty(user)
				ts.EqualError(sql.ErrNoRows, err.Error())
			},
		},
	}

	for _, c := range cases {
		ts.Run(c.name, func() {
			c.mockRepoBehaviour(ts.mock)
		})
	}
}
