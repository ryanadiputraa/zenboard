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
	"github.com/stretchr/testify/assert"
)

func randomUser() domain.User {
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

func TestCreateUser(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		assert.FailNow(t, "fail to create mock db: %s", err)
	}
	defer conn.Close()

	mockDB := db.New(conn)
	repo := NewUserRepository(mockDB)

	cases := []struct {
		name              string
		mockRepoBehaviour func(t *testing.T, mock sqlmock.Sqlmock)
	}{
		{
			name: "should create a user",
			mockRepoBehaviour: func(t *testing.T, mock sqlmock.Sqlmock) {
				user := randomUser()
				assert := assert.New(t)

				mock.ExpectQuery("INSERT INTO users").WillReturnRows(
					sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "picture", "locale", "board_limit", "created_at"}).
						AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.Picture, user.Locale, user.BoardLimit, user.CreatedAt),
				)

				created, err := repo.Save(context.TODO(), user)
				assert.NoError(err)
				assert.NotEmpty(created)
			},
		},
		{
			name: "should fail to create user",
			mockRepoBehaviour: func(t *testing.T, mock sqlmock.Sqlmock) {
				user := randomUser()
				assert := assert.New(t)

				mock.ExpectQuery("INSERT INTO users").WillReturnError(sql.ErrNoRows)

				created, err := repo.Save(context.TODO(), user)
				assert.EqualError(err, sql.ErrNoRows.Error())
				assert.Empty(created)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockRepoBehaviour(t, mock)
		})
	}
}
