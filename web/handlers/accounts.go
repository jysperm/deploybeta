package handlers

import (
  "strings"
  "github.com/kataras/iris"
)

import (
  . "github.com/jysperm/deploying/web/handlers/helpers"
  accountModel "github.com/jysperm/deploying/lib/models/account"
)

func RegisterAccount(ctx *iris.Context) {
  account := &accountModel.Account{}

  err := ctx.ReadJSON(account)

  if err != nil {
    ctx.JSON(iris.StatusBadRequest, NewHttpError(err))
    return
  }

  err = accountModel.Register(account)

  if err != nil && strings.Contains(err.Error(), "Key already exists") {
    ctx.JSON(iris.StatusConflict, NewHttpError(err))
    return
  }

  ctx.JSON(iris.StatusCreated, NewAccountResponse(account))
}
