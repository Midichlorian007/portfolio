package repository

import (
	"database/sql"
	"errors"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"main/internal/model"
)

func TestGetUserDb_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := &Repo{repoClient: db}
	expectedUser := &model.User{Id: 1, Name: "MinIO"}

	mock.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(expectedUser.Id, expectedUser.Name))

	user, err := repository.GetUserDb(1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.Id, user.Id)
	assert.Equal(t, expectedUser.Name, user.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserDb_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := &Repo{repoClient: db}

	mock.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	user, err := repository.GetUserDb(1)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "GetUserDb: user not found with id: "+strconv.Itoa(1))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserDb_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := &Repo{repoClient: db}

	mock.ExpectQuery("SELECT id, name FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnError(errors.New("database error"))

	user, err := repository.GetUserDb(1)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "GetUserDb: repo.repoClient: database error")

	assert.NoError(t, mock.ExpectationsWereMet())
}
