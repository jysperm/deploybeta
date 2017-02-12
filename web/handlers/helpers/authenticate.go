package helpers

import (
	"net/http"

	"github.com/labstack/echo"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	sessionModel "github.com/jysperm/deploying/lib/models/session"
)

func AuthenticateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		session, err := sessionModel.FindByToken(ctx.Request().Header.Get("Authorization"))

		if err != nil {
			return NewHTTPError(http.StatusUnauthorized, err)
		}

		account, err := accountModel.FindByName(session.Username)

		if err != nil {
			return NewHTTPError(http.StatusUnauthorized, err)
		}

		ctx.Set("account", account)

		return next(ctx)
	}
}

func GetSessionAccount(ctx echo.Context) *accountModel.Account {
	return ctx.Get("account").(*accountModel.Account)
}
