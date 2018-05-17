package handlers

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/jysperm/deploybeta/lib/models"
	. "github.com/jysperm/deploybeta/web/handlers/helpers"
)

func CreateSession(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	account, err := models.FindAccountByName(params["username"])

	if err != nil && err == models.ErrAccountNotFound {
		return NewHTTPError(http.StatusUnauthorized, err)
	} else if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	err = account.ComparePassword(params["password"])

	if err != nil {
		return NewHTTPError(http.StatusUnauthorized, err)
	}

	session, err := models.CreateSession(&account)

	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewSessionResponse(session))
}
