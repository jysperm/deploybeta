package handlers

import (
	"net/http"

	"github.com/labstack/echo"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	. "github.com/jysperm/deploying/web/handlers/helpers"
)

func RegisterAccount(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	account := &accountModel.Account{
		Username: params["username"],
		Email:    params["email"],
	}

	err = accountModel.Register(account, params["password"])

	if err != nil && err == accountModel.ErrUsernameConflict {
		return NewHTTPError(http.StatusConflict, err)
	} else if err != nil && err == accountModel.ErrInvalidUsername {
		return NewHTTPError(http.StatusBadRequest, err)
	} else if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewAccountResponse(account))
}

func CurrentAccount(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, NewAccountResponse(GetSessionAccount(ctx)))
}
