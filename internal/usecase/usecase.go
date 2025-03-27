package usecase

import (
	"errors"
	"io"
	"main/internal/model"
	"main/internal/repository"
)

func NewUseCase(cfg *model.Config) (UseCase, []io.Closer, error) {
	useCases := useCase{
		cfg: cfg,
	}
	repo, closers, err := repository.NewRepo(cfg)
	if err != nil {
		return &useCases, closers, err
	}

	useCases.repository = repo

	return &useCases, closers, nil
}

type UseCase interface {
	CreateUser(user *model.User) error
	GetUser(id int) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id int) error
}
type useCase struct {
	cfg        *model.Config
	repository *repository.Repo
}

func (uc *useCase) CreateUser(user *model.User) error {

	if user == nil {
		return errors.New(model.LevelError + "CreateUserDb: user is nil")
	}

	if user.Name == "" {
		return errors.New(model.LevelError + "CreateUserDb: user name is empty")
	}

	return uc.repository.CreateUserDb(user)
}

func (uc *useCase) GetUser(id int) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New(model.LevelError + "GetUser: user id invalid")
	}
	return uc.repository.GetUserDb(id)
}

func (uc *useCase) UpdateUser(user *model.User) error {
	if user == nil {
		return errors.New(model.LevelError + "UpdateUser: user is nil")
	}

	if user.Id == 0 {
		return errors.New(model.LevelError + "UpdateUser: user id is invalid")
	}

	if user.Name == "" {
		return errors.New(model.LevelError + "UpdateUser: user name is empty")
	}

	return uc.repository.UpdateUserDb(user)
}

func (uc *useCase) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New(model.LevelError + "DeleteUser: user id invalid")
	}
	return uc.repository.DeleteUserDb(id)
}
