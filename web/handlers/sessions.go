package handlers

import (
	"github.com/kataras/iris"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	sessionModel "github.com/jysperm/deploying/lib/models/session"
	. "github.com/jysperm/deploying/web/handlers/helpers"
)

func CreateSession(ctx *iris.Context) {
	params := map[string]string{}
	err := ctx.ReadJSON(&params)

	if err != nil {
		ctx.JSON(iris.StatusBadRequest, NewHttpError(err))
		return
	}

	account, err := accountModel.FindByName(params["username"])

	if err != nil && err == accountModel.ErrAccountNotFound {
		ctx.JSON(iris.StatusUnauthorized, NewHttpError(err))
		return
	} else if err != nil {
		ctx.JSON(iris.StatusInternalServerError, NewHttpError(err))
		return
	}

	err = account.ComparePassword(params["password"])

	if err != nil {
		ctx.JSON(iris.StatusUnauthorized, NewHttpError(err))
		return
	}

	session, err := sessionModel.CreateToken(account)

	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, NewHttpError(err))
		return
	}

	ctx.JSON(iris.StatusCreated, NewSessionResponse(session))
}
