package handlers

import (
	"net/http"

	"github.com/labstack/echo"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	sessionModel "github.com/jysperm/deploying/lib/models/session"
	. "github.com/jysperm/deploying/web/handlers/helpers"
)

func CreateSession(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	account, err := accountModel.FindByName(params["username"])

	if err != nil && err == accountModel.ErrAccountNotFound {
		return NewHTTPError(http.StatusUnauthorized, err)
	} else if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	err = account.ComparePassword(params["password"])

	if err != nil {
		return NewHTTPError(http.StatusUnauthorized, err)
	}

	session, err := sessionModel.CreateToken(account)

	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewSessionResponse(session))
}
