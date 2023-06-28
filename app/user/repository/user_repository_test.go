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

func TestSave(t *testing.T) {
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
				user := createRandomUser()
				assert := assert.New(t)

				mock.ExpectQuery("INSERT INTO users").WillReturnRows(
					sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "picture", "locale", "board_limit", "created_at"}).
						AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.Picture, user.Locale, user.BoardLimit, user.CreatedAt),
				)

				created, err := repo.Save(context.TODO(), user)
				assert.NoError(err)
				assert.NotEmpty(created)

				assert.Equal(user.ID, created.ID)
				assert.Equal(user.FirstName, created.FirstName)
				assert.Equal(user.LastName, created.LastName)
				assert.Equal(user.Email, created.Email)
				assert.Equal(user.Picture, created.Picture)
				assert.Equal(user.Locale, created.Locale)
				assert.Equal(user.BoardLimit, created.BoardLimit)
				assert.NotZero(created.CreatedAt)
			},
		},
		{
			name: "should fail to create user",
			mockRepoBehaviour: func(t *testing.T, mock sqlmock.Sqlmock) {
				user := createRandomUser()
				assert := assert.New(t)

				mock.ExpectQuery("INSERT INTO users").WillReturnError(sql.ErrNoRows)

				created, err := repo.Save(context.TODO(), user)
				assert.EqualError(sql.ErrNoRows, err.Error())
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

func TestList(t *testing.T) {
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
			name: "should return list of users within array of id(s)",
			mockRepoBehaviour: func(t *testing.T, mock sqlmock.Sqlmock) {
				var users []domain.User
				var ids []string
				assert := assert.New(t)

				for i := 0; i < 5; i++ {
					user := createRandomUser()
					users = append(users, user)
					ids = append(ids, user.ID)
				}

				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ANY").WillReturnRows(
					sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "picture", "locale", "board_limit", "created_at"}).
						AddRow(users[0].ID, users[0].FirstName, users[0].LastName, users[0].Email, users[0].Picture, users[0].Locale, users[0].BoardLimit, users[0].CreatedAt).
						AddRow(users[1].ID, users[1].FirstName, users[1].LastName, users[1].Email, users[1].Picture, users[1].Locale, users[1].BoardLimit, users[1].CreatedAt).
						AddRow(users[2].ID, users[2].FirstName, users[2].LastName, users[2].Email, users[2].Picture, users[2].Locale, users[2].BoardLimit, users[2].CreatedAt).
						AddRow(users[3].ID, users[3].FirstName, users[3].LastName, users[3].Email, users[3].Picture, users[3].Locale, users[3].BoardLimit, users[3].CreatedAt).
						AddRow(users[4].ID, users[4].FirstName, users[4].LastName, users[4].Email, users[4].Picture, users[4].Locale, users[4].BoardLimit, users[4].CreatedAt),
				)

				list, err := repo.List(context.TODO(), ids)
				assert.NoError(err)
				assert.Len(list, 5)

				for i, v := range list {
					assert.NotEmpty(v)
					assert.Equal(users[i].ID, v.ID)
					assert.Equal(users[i].FirstName, v.FirstName)
					assert.Equal(users[i].LastName, v.LastName)
					assert.Equal(users[i].Email, v.Email)
					assert.Equal(users[i].Picture, v.Picture)
					assert.Equal(users[i].Locale, v.Locale)
					assert.Equal(users[i].BoardLimit, v.BoardLimit)
					assert.NotZero(v.CreatedAt)
				}
			},
		},
		{
			name: "should fail to list users that didn't exists",
			mockRepoBehaviour: func(t *testing.T, mock sqlmock.Sqlmock) {
				var ids []string
				assert := assert.New(t)

				for i := 0; i < 5; i++ {
					user := createRandomUser()
					ids = append(ids, user.ID)
				}
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ANY").WillReturnError(sql.ErrNoRows)

				users, err := repo.List(context.TODO(), ids)
				assert.Empty(users)
				assert.EqualError(sql.ErrNoRows, err.Error())
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockRepoBehaviour(t, mock)
		})
	}
}
func TestFindByID(t *testing.T) {
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
			name: "should return a user data with correct id",
			mockRepoBehaviour: func(t *testing.T, mock sqlmock.Sqlmock) {
				user := createRandomUser()
				assert := assert.New(t)

				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(
					sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "picture", "locale", "board_limit", "created_at"}).
						AddRow(user.ID, user.FirstName, user.LastName, user.Email, user.Picture, user.Locale, user.BoardLimit, user.CreatedAt),
				)

				data, err := repo.FindByID(context.TODO(), user.ID)
				assert.NotEmpty(data)
				assert.NoError(err)

				assert.Equal(user.ID, data.ID)
				assert.Equal(user.FirstName, data.FirstName)
				assert.Equal(user.LastName, data.LastName)
				assert.Equal(user.Email, data.Email)
				assert.Equal(user.Picture, data.Picture)
				assert.Equal(user.Locale, data.Locale)
				assert.Equal(user.BoardLimit, data.BoardLimit)
				assert.NotZero(data.CreatedAt)
			},
		},
		{
			name: "should fail to find user that didn't exists",
			mockRepoBehaviour: func(t *testing.T, mock sqlmock.Sqlmock) {
				assert := assert.New(t)
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnError(sql.ErrNoRows)

				user, err := repo.FindByID(context.TODO(), uuid.NewString())
				assert.Empty(user)
				assert.EqualError(sql.ErrNoRows, err.Error())
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.mockRepoBehaviour(t, mock)
		})
	}
}
