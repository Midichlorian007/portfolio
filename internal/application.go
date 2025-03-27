package internal

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	"main/internal/handler"
	"main/internal/middleware"
	"main/internal/model"
	"os"
)

type Application interface {
	Shutdown() error
	Start() error
}

type application struct {
	server      *echo.Echo
	handlers    *handler.Handler
	middlewares *middleware.Middleware
	cfg         *model.Config
}

func NewApplication() Application {
	return &application{
		server: echo.New(),
		cfg:    &model.Config{},
	}
}

func (app *application) Start() error {
	err := app.parseConfig(model.ConfigPath)
	if err != nil {
		return err
	}

	app.initHandlers()
	app.initMiddlewares()
	app.initRouts()

	err = app.server.Start(app.cfg.Server.Port)
	if err != nil {
		return errors.New(model.LevelError + "server.Start: " + err.Error())
	}

	return nil
}

func (app *application) Shutdown() error {
	if app.server == nil {
		return errors.New(model.LevelError + "Shutdown: server is nil ")
	}

	if err := app.server.Shutdown(context.TODO()); err != nil {
		return errors.New(model.LevelError + "Shutdown: server.Shutdown: " + err.Error())
	}

	return nil
}

func (app *application) initRouts() {

	app.server.GET("/live", app.handlers.Live)
	app.server.GET("/cfg", app.handlers.GetConfig)
	app.server.GET("/sleep/:seconds", app.handlers.Sleep)

}

func (app *application) initMiddlewares() {
	app.middlewares = middleware.NewMiddleware(app.cfg)

	app.server.Use(echoMw.RequestIDWithConfig(echoMw.RequestIDConfig{
		Generator: func() string {
			return uuid.NewString()
		},
	}))

	app.server.Use(echoMw.BodyDumpWithConfig(echoMw.BodyDumpConfig{
		Handler: app.middlewares.BodyLogger,
	}))

	//app.server.Use(app.middlewares.Authorization)
	app.server.Use(echoMw.Recover())
}

func (app *application) initHandlers() {
	app.handlers = handler.NewHandler(app.cfg)
}

func (app *application) parseConfig(path string) error {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return errors.New(model.LevelError + "parseConfig: os.ReadFile: " + err.Error())
	}

	err = json.Unmarshal(fileBytes, &app.cfg)
	if err != nil {
		return errors.New(model.LevelError + "parseConfig: json.Unmarshal: " + err.Error())
	}

	return nil
}
