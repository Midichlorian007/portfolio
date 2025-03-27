package internal

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	"io"
	"log"
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
	server *echo.Echo
	cfg    *model.Config

	handlers    *handler.Handler
	middlewares *middleware.Middleware

	AllClosers []io.Closer
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

	defer app.CloseAllClosers()

	err = app.initHandlers()
	if err != nil {
		return err
	}

	app.initRoutsAndMiddlewares()

	return app.server.Start(app.cfg.Server.Port)
}

func (app *application) CloseAllClosers() {
	if app.AllClosers == nil {
		return
	}

	for _, closer := range app.AllClosers {
		if err := closer.Close(); err != nil {
			log.Println(model.LevelError + "CloseAllClosers: " + err.Error())
		}
	}
}

func (app *application) Shutdown() error {
	if app.server != nil {
		if err := app.server.Shutdown(context.TODO()); err != nil {
			return errors.New(model.LevelError + "Shutdown: server.Shutdown: " + err.Error())
		}
	}

	return nil
}

func (app *application) initRoutsAndMiddlewares() {
	app.initBasicMiddlewares()

	//app.server.Use(app.middlewares.Authorization)
	app.server.GET("/live", app.handlers.Live)
	app.server.POST("/users", app.handlers.CreateUser)
	app.server.GET("/users/:user_id", app.handlers.GetUserById)
	app.server.PUT("/users/:user_id", app.handlers.UpdateUserById)
	app.server.DELETE("/users/:user_id", app.handlers.DeleteUserById)

	adminGroup := app.server.Group("/admin")
	//adminGroup.Use(app.middlewares.AdminAuthorization)
	adminGroup.GET("/cfg", app.handlers.GetConfig)
	adminGroup.GET("/sleep/:seconds", app.handlers.Sleep)
}

func (app *application) initBasicMiddlewares() {
	app.middlewares = middleware.NewMiddleware(app.cfg)

	app.server.Use(echoMw.RequestIDWithConfig(echoMw.RequestIDConfig{
		Generator: func() string {
			return uuid.NewString()
		},
	}))

	app.server.Use(echoMw.BodyDumpWithConfig(echoMw.BodyDumpConfig{
		Handler: app.middlewares.BodyLogger,
	}))

	app.server.Use(echoMw.Recover())
}

func (app *application) initHandlers() error {
	handlers, closers, err := handler.NewHandler(app.cfg)
	app.addClosers(closers)
	if err != nil {
		return err
	}

	app.handlers = handlers

	return nil
}

func (app *application) addClosers(closers []io.Closer) {

	if app.AllClosers == nil {
		app.AllClosers = make([]io.Closer, 0)
	}

	if closers == nil {
		return
	}

	app.AllClosers = append(app.AllClosers, closers...)

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
