package handler

import (
	"github.com/labstack/echo/v4"
	"io"
	"main/internal/model"
	"main/internal/usecase"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	cfg *model.Config

	useCase usecase.UseCase
}

func NewHandler(cfg *model.Config) (*Handler, []io.Closer, error) {
	handlers := Handler{
		cfg: cfg,
	}

	useCase, closers, err := usecase.NewUseCase(cfg)
	if err != nil {
		return &handlers, closers, err
	}

	handlers.useCase = useCase

	return &handlers, closers, nil
}

func (h *Handler) Live(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, *h.cfg)
}

func (h *Handler) Sleep(c echo.Context) error {
	secondsStr := c.Param("seconds")

	seconds, err := strconv.Atoi(secondsStr)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	time.Sleep(time.Duration(seconds) * time.Second)

	return c.NoContent(http.StatusOK)
}

func (h *Handler) CreateUser(c echo.Context) error {
	user := &model.User{}

	err := c.Bind(user)
	if err != nil {
		return c.String(http.StatusBadRequest, model.LevelError+"CreateUser: c.Bind: "+err.Error())
	}

	if user.Id != 0 {
		return c.String(http.StatusBadRequest, model.LevelError+"CreateUser: user id must not be presented")
	}

	err = h.useCase.CreateUser(user)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) UpdateUserById(c echo.Context) error {

	user := &model.User{}

	err := c.Bind(user)
	if err != nil {
		return c.String(http.StatusBadRequest, model.LevelError+"UpdateUserById: c.Bind: "+err.Error())
	}

	userIdStr := c.Param("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId == 0 {
		return c.String(http.StatusBadRequest, model.LevelError+"UpdateUserById: Invalid user id")
	}

	user.Id = userId

	err = h.useCase.UpdateUser(user)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) DeleteUserById(c echo.Context) error {
	userIdStr := c.Param("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId == 0 {
		return c.String(http.StatusBadRequest, model.LevelError+"DeleteUserById: Invalid user id")
	}

	err = h.useCase.DeleteUser(userId)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUserById(c echo.Context) error {
	userIdStr := c.Param("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return c.String(http.StatusBadRequest, model.LevelError+"GetUserById: Invalid user id")
	}

	user, err := h.useCase.GetUser(userId)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}
