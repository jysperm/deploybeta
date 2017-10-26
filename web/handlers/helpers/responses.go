package helpers

import "github.com/jysperm/deploying/lib/models"

type ErrorResponse struct {
	Error string `json:"error"`
}

type AccountResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AppResponse struct {
	Name          string           `json:"name"`
	Owner         string           `json:"owner"`
	GitRepository string           `json:"gitRepository"`
	Instances     int              `json:"instances"`
	Version       string           `json:"version"`
	Versions      []models.Version `json:"versions"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error: err.Error(),
	}
}

func NewAccountResponse(account *models.Account) AccountResponse {
	return AccountResponse{
		Username: account.Username,
		Email:    account.Email,
	}
}

func NewSessionResponse(session *models.Session) models.Session {
	return *session
}

func NewAppResponse(app *models.Application) AppResponse {
	appRes := AppResponse{
		Name:          app.Name,
		Owner:         app.Owner,
		Version:       app.Version,
		GitRepository: app.GitRepository,
		Instances:     app.Instances,
	}

	versions, err := models.ListVersions(app)
	if err != nil {
		return AppResponse{}
	}
	appRes.Versions = *versions

	return appRes
}

func NewAppsResponse(apps []models.Application) []AppResponse {
	appsRes := make([]AppResponse, 0)
	var app AppResponse
	for _, v := range apps {
		app.GitRepository = v.GitRepository
		app.Owner = v.Owner
		app.Name = v.Name
		app.Version = v.Version
		app.Instances = v.Instances
		versions, err := models.ListVersions(&v)
		if err != nil {
			panic(err)
		}
		app.Versions = *versions
		appsRes = append(appsRes, app)
	}

	return appsRes
}

func NewVersionResponse(version *models.Version) models.Version {
	return *version
}
