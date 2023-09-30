package app

import (
	"os"

	v1 "github.com/evermake/git-diff-view/internal/controller/http/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var _defaultAddr = ":7777"

type Option func(*App)

func WithAddr(addr string) Option {
	return func(app *App) {
		app.addr = addr
	}
}

func WithRepoPath(repoPath string) Option {
	return func(app *App) {
		app.repoPath = repoPath
	}
}

type App struct {
	addr     string
	repoPath string
}

func New(options ...Option) (*App, error) {
	app := &App{addr: _defaultAddr}

	for _, option := range options {
		option(app)
	}

	if app.repoPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		app.repoPath = wd
	}

	return app, nil
}

func (a *App) Run() error {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	if err := v1.RegisterHandlers(e, v1.NewServer(a.repoPath)); err != nil {
		return err
	}

	return e.Start(a.addr)
}
