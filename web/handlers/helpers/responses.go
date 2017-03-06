package helpers

import (
	accountModel "github.com/jysperm/deploying/lib/models/account"
	appModel "github.com/jysperm/deploying/lib/models/app"
	sessionModel "github.com/jysperm/deploying/lib/models/session"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type AccountResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error: err.Error(),
	}
}

func NewAccountResponse(account *accountModel.Account) AccountResponse {
	return AccountResponse{
		Username: account.Username,
		Email:    account.Email,
	}
}

func NewSessionResponse(session *sessionModel.Session) sessionModel.Session {
	return *session
}

func NewAppResponse(app *appModel.Application) appModel.Application {
	return *app
}

type ImageResponse struct {
	Shasum  string
	Version string
}

func NewImageResponse(shasum string, version string) ImageResponse {
	return ImageResponse{
		Shasum:  shasum,
		Version: version,
	}
}
