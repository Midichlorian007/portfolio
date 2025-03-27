package handler

import (
	"github.com/labstack/echo/v4"
	"main/internal/model"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	cfg *model.Config
}

func NewHandler(cfg *model.Config) *Handler {
	return &Handler{cfg: cfg}
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
