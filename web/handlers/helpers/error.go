package helpers

import (
	"net/http"

	"github.com/labstack/echo"
)

func NewHTTPError(code int, err error) error {
	return echo.NewHTTPError(http.StatusConflict, NewErrorResponse(err))
}

func HTTPErrorHandler(err error, ctx echo.Context) {
	var msg interface{}

	code := http.StatusInternalServerError

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else if ctx.Echo().Debug {
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}
	if _, ok := msg.(string); ok {
		msg = map[string]interface{}{"message": msg}
	}

	ctx.Echo().Logger.Error(ctx.Path(), " ", err)

	if !ctx.Response().Committed {
		if ctx.Request().Method == echo.HEAD {
			err = ctx.NoContent(code)
		} else {
			err = ctx.JSON(code, msg)
		}
		if err != nil {
			ctx.Echo().Logger.Error(ctx.Path(), " ", err)
		}
	}
}
