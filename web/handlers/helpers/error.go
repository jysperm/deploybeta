package helpers

import (
	"net/http"

	"github.com/labstack/echo"
)

func NewHTTPError(code int, err error) error {
	return echo.NewHTTPError(http.StatusConflict, NewErrorResponse(err))
}
