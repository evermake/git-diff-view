package app

import (
	v1 "github.com/evermake/git-diff-view/internal/controller/http/v1"
	"github.com/labstack/echo/v4"
)

var _defaultAddr = ":7777"

type Config struct {
	addr string
}

type Option func(config *Config)

func WithAddr(addr string) Option {
	return func(config *Config) {
		config.addr = addr
	}
}

func Run(options ...Option) error {
	config := Config{
		addr: _defaultAddr,
	}

	for _, option := range options {
		option(&config)
	}

	e := echo.New()

	if err := v1.RegisterHandlers(e); err != nil {
		return err
	}

	return e.Start(config.addr)
}
