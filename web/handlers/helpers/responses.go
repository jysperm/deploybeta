package helpers

import (
	"fmt"

	"github.com/jysperm/deploying/lib/models"
	"github.com/jysperm/deploying/lib/swarm"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type AccountResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AppResponse struct {
	Name          string            `json:"name"`
	Owner         string            `json:"owner"`
	GitRepository string            `json:"gitRepository"`
	Instances     int               `json:"instances"`
	Version       string            `json:"version"`
	Versions      []models.Version  `json:"versions"`
	Nodes         []swarm.Container `json:"nodes"`
}

type DataSourceResponse struct {
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	Type      string `json:"type"`
	Instances int    `json:"instances"`
}

type DataSourceNodeResponse struct {
	Host string `json:"host"`
	Role string `json:"role"`
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

	versions, _ := models.ListVersions(app)
	appRes.Versions = *versions

	return appRes
}

func NewDataSourceResponse(dataSource *models.DataSource) DataSourceResponse {
	return DataSourceResponse{
		Name:      dataSource.Name,
		Owner:     dataSource.Owner,
		Type:      dataSource.Type,
		Instances: dataSource.Instances,
	}
}

func NewDataSourceNodeResponse(dataSource *models.DataSourceNode) DataSourceNodeResponse {
	return DataSourceNodeResponse{
		Host: dataSource.Host,
		Role: dataSource.Role,
	}
}

func NewDataSourcesResponse(dataSources []models.DataSource) []DataSourceResponse {
	result := make([]DataSourceResponse, 0)

	for _, dataSource := range dataSources {
		result = append(result, NewDataSourceResponse(&dataSource))
	}

	return result
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
			fmt.Println(err)
		}
		app.Versions = *versions
		nodes, err := swarm.ListContainers(&v)
		if err != nil {
			fmt.Println(err.Error())
		}
		if nodes == nil {
			app.Nodes = []swarm.Container{}
		} else {
			app.Nodes = *nodes
		}
		appsRes = append(appsRes, app)
	}

	return appsRes
}

func NewVersionResponse(version *models.Version) models.Version {
	return *version
}
