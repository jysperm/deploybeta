package helpers

import (
	"github.com/kataras/iris"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	sessionModel "github.com/jysperm/deploying/lib/models/session"
)

type AuthenticateMiddleware struct{}

func (middleware *AuthenticateMiddleware) Serve(ctx *iris.Context) {
	session, err := sessionModel.FindByToken(ctx.RequestHeader("Authorization"))

	if err != nil {
		ctx.JSON(iris.StatusUnauthorized, NewHttpError(err))
		return
	}

	account, err := accountModel.FindByName(session.Username)

	if err != nil {
		ctx.JSON(iris.StatusUnauthorized, NewHttpError(err))
		return
	}

	ctx.Set("account", account)
	ctx.Next()
}
