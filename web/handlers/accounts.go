package handlers

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/jysperm/deploybeta/lib/models"
	. "github.com/jysperm/deploybeta/web/handlers/helpers"
)

func RegisterAccount(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	account := &models.Account{
		Username: params["username"],
		Email:    params["email"],
	}

	err = models.RegisterAccount(account, params["password"])

	if err != nil && err == models.ErrUsernameConflict {
		return NewHTTPError(http.StatusConflict, err)
	} else if err != nil && err == models.ErrInvalidUsername {
		return NewHTTPError(http.StatusBadRequest, err)
	} else if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewAccountResponse(account))
}

func CurrentAccount(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, NewAccountResponse(GetSessionAccount(ctx)))
}
