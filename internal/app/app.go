package app

import (
	v1 "github.com/evermake/git-diff-view/internal/controller/http/v1"
	"github.com/labstack/echo/v4"
)

var _defaultAddr = ":7777"

type Option func(*App)

func WithAddr(addr string) Option {
	return func(app *App) {
		app.addr = addr
	}
}

type App struct {
	addr string
}

func New(options ...Option) App {
	app := App{addr: _defaultAddr}

	for _, option := range options {
		option(&app)
	}

	return app
}

func (a App) Run() error {
	e := echo.New()

	if err := v1.RegisterHandlers(e); err != nil {
		return err
	}

	return e.Start(a.addr)
}
