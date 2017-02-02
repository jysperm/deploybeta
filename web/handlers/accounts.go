package handlers

import (
  "github.com/kataras/iris"
)

import . "github.com/jysperm/deploying/lib/models"

func RegisterAccount(ctx *iris.Context) {
  account := &Account{}

  err := ctx.ReadJSON(account)

  if err != nil {
		panic(err)
	}

  CreateAccount(account)

  ctx.JSON(iris.StatusCreated, account)
}
