package helpers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	appModel "github.com/jysperm/deploying/lib/models/app"
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

func AppOwnerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		account := GetSessionAccount(ctx)
		appName := ctx.Param("name")
		apps, err := appModel.GetAppsOfAccount(account)

		if err != nil {
			return NewHTTPError(http.StatusInternalServerError, err)
		}

		if err == nil && len(apps) == 0 {
			return NewHTTPError(http.StatusUnauthorized, errors.New("Not found application"))
		}

		for _, app := range apps {
			if app.Name == appName {
				ctx.Set("app", app)
				return next(ctx)
			}
		}

		return NewHTTPError(http.StatusBadRequest, errors.New("Not found application"))
	}
}
func GetSessionAccount(ctx echo.Context) *accountModel.Account {
	return ctx.Get("account").(*accountModel.Account)
}
