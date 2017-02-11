package handlers

import (
	"github.com/kataras/iris"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	. "github.com/jysperm/deploying/web/handlers/helpers"
)

func RegisterAccount(ctx *iris.Context) {
	params := map[string]string{}
	err := ctx.ReadJSON(&params)

	if err != nil {
		ctx.JSON(iris.StatusBadRequest, NewHttpError(err))
		return
	}

	account := &accountModel.Account{
		Username: params["username"],
		Email:    params["email"],
	}

	err = accountModel.Register(account, params["password"])

	if err != nil && err == accountModel.ErrUsernameConflict {
		ctx.JSON(iris.StatusConflict, NewHttpError(err))
		return
	} else if err != nil && err == accountModel.ErrInvalidUsername {
		ctx.JSON(iris.StatusBadRequest, NewHttpError(err))
		return
	} else if err != nil {
		ctx.JSON(iris.StatusInternalServerError, NewHttpError(err))
		return
	}

	ctx.JSON(iris.StatusCreated, NewAccountResponse(account))
}

func CurrentAccount(ctx *iris.Context) {
	account := ctx.Get("account").(*accountModel.Account)

	ctx.JSON(iris.StatusOK, NewAccountResponse(account))
}
