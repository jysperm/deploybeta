package helpers

import (
	accountModel "github.com/jysperm/deploying/lib/models/account"
	appModel "github.com/jysperm/deploying/lib/models/app"
	sessionModel "github.com/jysperm/deploying/lib/models/session"
	versionModel "github.com/jysperm/deploying/lib/models/version"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type AccountResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AppResponse struct {
	Name          string                 `json:"name"`
	Owner         string                 `json:"owner"`
	GitRepository string                 `json:"gitRepository"`
	Instances     int                    `json:"instances"`
	Version       string                 `json:"version"`
	Versions      []versionModel.Version `json:"versions"`
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

func NewAppResponse(app *appModel.Application) AppResponse {
	appRes := AppResponse{
		Name:          app.Name,
		Owner:         app.Owner,
		Version:       app.Version,
		GitRepository: app.GitRepository,
		Instances:     app.Instances,
	}

	versions, err := versionModel.ListAll(*app)
	if err != nil {
		return AppResponse{}
	}
	appRes.Versions = *versions

	return appRes
}

func NewAppsResponse(apps []appModel.Application) []AppResponse {
	var appsRes []AppResponse
	var app AppResponse
	for _, v := range apps {
		app.GitRepository = v.GitRepository
		app.Owner = v.Owner
		app.Name = v.Name
		app.Version = v.Version
		app.Instances = v.Instances
		versions, err := versionModel.ListAll(v)
		if err != nil {
			panic(err)
		}
		app.Versions = *versions
		appsRes = append(appsRes, app)
	}

	return appsRes
}

func NewVersionResponse(version *versionModel.Version) versionModel.Version {
	return *version
}
