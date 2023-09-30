package v1

import (
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/evermake/git-diff-view/internal/controller/http/v1/openapi"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) error {
	swagger, err := openapi.GetSwagger()
	if err != nil {
		return err
	}
	swagger.Servers = nil

	e.Use(middleware.OapiRequestValidator(swagger))

	openapi.RegisterHandlers(e, openapi.NewStrictHandler(
		NewServer(),
		nil,
	))

	return nil
}
