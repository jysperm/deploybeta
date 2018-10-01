package helpers

import (
	"github.com/jysperm/deploybeta/lib/models"
	"github.com/jysperm/deploybeta/lib/swarm"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type AccountResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ApplicationResponse struct {
	Name          string            `json:"name"`
	OwnerUsername string            `json:"ownerUsername"`
	GitRepository string            `json:"gitRepository"`
	Instances     int               `json:"instances"`
	VersionTag    string            `json:"versionTag"`
	Versions      []models.Version  `json:"versions"`
	Nodes         []swarm.Container `json:"nodes"`
}

type DataSourceResponse struct {
	Name          string                `json:"name"`
	OwnerUsername string                `json:"ownerUsername"`
	Type          string                `json:"type"`
	Instances     int                   `json:"instances"`
	LinkedApps    []ApplicationResponse `json:"linkedApps"`
}

type DataSourceNodeResponse struct {
	Host       string `json:"host"`
	Role       string `json:"role"`
	MasterHost string `json:"masterHost"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		Error: err.Error(),
	}
}

func NewAccountResponse(account *models.Account) *AccountResponse {
	return &AccountResponse{
		Username: account.Username,
		Email:    account.Email,
	}
}

func NewApplicationResponse(app *models.Application, versions []models.Version, nodes []swarm.Container) *ApplicationResponse {
	return &ApplicationResponse{
		Name:          app.Name,
		OwnerUsername: app.OwnerUsername,
		GitRepository: app.GitRepository,
		Instances:     app.Instances,
		VersionTag:    app.VersionTag,
		Versions:      versions,
		Nodes:         nodes,
	}
}

func NewDataSourceResponse(dataSource *models.DataSource, apps []models.Application) *DataSourceResponse {
	appsResponses := make([]ApplicationResponse, len(apps))

	for i, app := range apps {
		appsResponses[i] = *NewApplicationResponse(&app, nil, nil)
	}

	return &DataSourceResponse{
		Name:          dataSource.Name,
		OwnerUsername: dataSource.OwnerUsername,
		Type:          dataSource.Type,
		Instances:     dataSource.Instances,
		LinkedApps:    appsResponses,
	}
}

func NewDataSourceNodeResponse(node *models.DataSourceNode, command *models.DataSourceNodeCommand) *DataSourceNodeResponse {
	nodeResponse := &DataSourceNodeResponse{
		Host:       node.Host,
		Role:       node.Role,
		MasterHost: node.MasterHost,
	}

	if command != nil && command.Command == models.COMMAND_CHANGE_ROLE {
		nodeResponse.Role = command.Role
		nodeResponse.MasterHost = command.MasterHost
	}

	return nodeResponse
}

func NewDataSourceNodesResponse(nodes []models.DataSourceNode) []DataSourceNodeResponse {
	responses := make([]DataSourceNodeResponse, len(nodes))

	for i, node := range nodes {
		responses[i] = *NewDataSourceNodeResponse(&node, nil)
	}

	return responses
}

func NewSessionResponse(session *models.Session) *models.Session {
	return session
}

func NewVersionResponse(version *models.Version) *models.Version {
	return version
}
