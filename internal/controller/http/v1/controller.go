package v1

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/evermake/git-diff-view/internal/controller/http/v1/openapi"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Group) error {
	swagger, err := openapi.GetSwagger()
	if err != nil {
		return err
	}
	swagger.Servers = nil

	{
		v1 := e.Group("/v1", middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
			ErrorHandler: func(e echo.Context, err *echo.HTTPError) error {
				var msg string
				if err.Code == http.StatusInternalServerError {
					// do not expose internal error message
					// as it can contain sensible data
					msg = "internal server error"
				} else {
					msg = fmt.Sprint(err.Message)
				}

				return e.JSON(err.Code, openapi.Error{
					Error: msg,
				})
			},
		}))

		openapi.RegisterHandlers(v1, openapi.NewStrictHandler(
			NewServer(),
			nil,
		))
	}

	return nil
}
