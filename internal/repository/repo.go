package repository

import (
	"errors"
	"io"
	"main/internal/model"
	"strconv"

	"database/sql"
	_ "modernc.org/sqlite"
)

type Repo struct {
	cfg        *model.Config
	repoClient *sql.DB
}

// NewRepo Намеренно использовал SQLite вместо Postgres, чтобы была автономность.
func NewRepo(cfg *model.Config) (*Repo, []io.Closer, error) {
	closers := make([]io.Closer, 0)
	repo := Repo{
		cfg: cfg,
	}

	db, err := sql.Open(cfg.Sqlite.Driver, cfg.Sqlite.Host)
	if err != nil {
		return &repo, closers, errors.New(model.LevelError + "NewRepo: sql.Open: " + err.Error())
	}
	repo.repoClient = db

	closers = append(closers)

	err = repo.CreateTable()
	if err != nil {
		return &repo, closers, err
	}

	return &repo, closers, nil
}

func (repo *Repo) CreateTable() error {
	_, err := repo.repoClient.Exec(createTableQry)

	if err != nil {
		return errors.New(model.LevelError + "CreateTable: repo.repoClient.Exec: " + err.Error())
	}

	return nil
}

func (repo *Repo) CreateUserDb(user *model.User) error {

	result, err := repo.repoClient.Exec(CreateUserQry, user.Name)
	if err != nil {
		return errors.New(model.LevelError + "CreateUserDb: repoClient.Exec: " + err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.New(model.LevelError + "CreateUserDb: result.LastInsertId: " + err.Error())
	}

	user.Id = int(id)

	return nil
}

func (repo *Repo) GetUserDb(id int) (*model.User, error) {
	user := &model.User{}

	err := repo.repoClient.QueryRow(GetUserQry, id).Scan(&user.Id, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(model.LevelError + "GetUserDb: user not found with id: " + strconv.Itoa(id))
		}
		return nil, errors.New(model.LevelError + "GetUserDb: repo.repoClient: " + err.Error())
	}

	return user, nil
}

func (repo *Repo) UpdateUserDb(user *model.User) error {

	result, err := repo.repoClient.Exec(UpdateUserQry, user.Name, user.Id)
	if err != nil {
		return errors.New(model.LevelError + "UpdateUserDb: repoClient.Exec: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New(model.LevelError + "UpdateUserDb: result.RowsAffected: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New(model.LevelError + "UpdateUserDb: user not found with id: " + strconv.Itoa(user.Id))
	}

	return nil
}

func (repo *Repo) DeleteUserDb(id int) error {

	result, err := repo.repoClient.Exec(DeleteUserQry, id)
	if err != nil {
		return errors.New(model.LevelError + "DeleteUserDb: repoClient.Exec: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New(model.LevelError + "DeleteUserDb: result.RowsAffected: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New(model.LevelError + "DeleteUserDb: user not found with id: " + strconv.Itoa(id))
	}

	return nil
}
