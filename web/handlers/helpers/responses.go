package helpers

import (
	"log"

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
	Name       string   `json:"name"`
	Owner      string   `json:"owner"`
	Type       string   `json:"type"`
	Instances  int      `json:"instances"`
	LinkedApps []string `json:"linkedApps"`
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

	nodes, _ := swarm.ListNodes(app)
	if len(nodes) == 0 {
		appRes.Nodes = []swarm.Container{}
	} else {
		appRes.Nodes = nodes
	}
	return appRes
}

func NewDataSourceResponse(dataSource *models.DataSource) DataSourceResponse {
	linkedApps, err := dataSource.GetLinkedAppNames()

	if err != nil {
		log.Println(err)
	}

	return DataSourceResponse{
		Name:       dataSource.Name,
		Owner:      dataSource.Owner,
		Type:       dataSource.Type,
		Instances:  dataSource.Instances,
		LinkedApps: linkedApps,
	}
}

func NewDataSourceNodeResponse(node *models.DataSourceNode) DataSourceNodeResponse {
	return DataSourceNodeResponse{
		Host: node.Host,
		Role: node.Role,
	}
}

func NewDataSourceNodesResponse(nodes []models.DataSourceNode) []DataSourceNodeResponse {
	result := make([]DataSourceNodeResponse, 0)

	for _, node := range nodes {
		result = append(result, NewDataSourceNodeResponse(&node))
	}

	return result
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
			log.Println(err)
		}
		app.Versions = *versions
		nodes, err := swarm.ListNodes(&v)
		if err != nil {
			log.Println(err.Error())
		}
		if len(nodes) == 0 {
			app.Nodes = []swarm.Container{}
		} else {
			app.Nodes = nodes
		}
		appsRes = append(appsRes, app)
	}

	return appsRes
}

func NewVersionResponse(version *models.Version) models.Version {
	return *version
}
