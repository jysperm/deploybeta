package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/buger/jsonparser"
	"github.com/hashicorp/errwrap"
	"github.com/labstack/echo"

	"github.com/jysperm/deploybeta/lib/models"
	"github.com/jysperm/deploybeta/lib/swarm"
	. "github.com/jysperm/deploybeta/web/handlers/helpers"
)

func ListDataSources(ctx echo.Context) error {
	account := GetSessionAccount(ctx)

	dataSources := make([]models.DataSource, 0)
	err := account.DataSources().FetchAll(&dataSources)

	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	responses := make([]DataSourceResponse, len(dataSources))

	for i, dataSource := range dataSources {
		apps := make([]models.Application, 0)
		err := dataSource.Apps().FetchAll(&apps)

		if err != nil {
			return err
		}

		responses[i] = *NewDataSourceResponse(&dataSource, apps)
	}

	return ctx.JSON(http.StatusOK, responses)
}

func CreateDataSource(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	dataSource := &models.DataSource{
		Name:          params["name"],
		Type:          params["type"],
		OwnerUsername: GetSessionAccount(ctx).Username,
		Instances:     2,
	}

	err = models.CreateDataSource(dataSource)

	if err != nil && err == models.ErrUpdateConflict {
		return NewHTTPError(http.StatusConflict, err)
	} else if err != nil && err == models.ErrInvalidName {
		return NewHTTPError(http.StatusBadRequest, err)
	} else if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := swarm.UpdateDataSource(dataSource); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewDataSourceResponse(dataSource, nil))
}

func UpdateDataSource(ctx echo.Context) error {
	dataSource := GetDataSourceInfo(ctx)
	jsonBuf := make([]byte, 1024)

	if _, err := ctx.Request().Body.Read(jsonBuf); err != nil && err != io.EOF {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	instances, valueType, _, err := jsonparser.Get(jsonBuf, "instances")
	if err != jsonparser.KeyPathNotFoundError && err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	if valueType != jsonparser.NotExist {
		realValue, err := strconv.Atoi(string(instances))
		if err != nil {
			return NewHTTPError(http.StatusBadRequest, err)
		}
		dataSource.Instances = realValue
	}

	if err := dataSource.UpdateInstances(dataSource.Instances); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := swarm.UpdateDataSource(dataSource); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, NewDataSourceResponse(dataSource, nil))
}

func LinkDataSource(ctx echo.Context) error {
	appName := ctx.Param("appName")
	dataSource := GetDataSourceInfo(ctx)

	app, err := models.FindAppByName(appName)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	if err := dataSource.LinkApp(app); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := swarm.UpdateAppService(app); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.String(http.StatusOK, "")
}

func UnlinkDataSource(ctx echo.Context) error {
	appName := ctx.Param("appName")
	dataSource := GetDataSourceInfo(ctx)

	app, err := models.FindAppByName(appName)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	if err := dataSource.UnlinkApp(app); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := swarm.UpdateAppService(app); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.String(http.StatusOK, "")
}

func DeleteDataSource(ctx echo.Context) error {
	dataSource := GetDataSourceInfo(ctx)

	if err := dataSource.Destroy(); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := swarm.RemoveDataSource(dataSource); err != nil {
		return NewHTTPError(http.StatusInternalServerError, errwrap.Wrapf("remove dataSource service: {{err}}", err))
	}

	return ctx.String(http.StatusOK, "")
}

func ListDataSourceNodes(ctx echo.Context) error {
	dataSource := GetDataSourceInfo(ctx)

	nodes := make([]models.DataSourceNode, 0)
	err := dataSource.Nodes().FetchAll(&nodes)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, NewDataSourceNodesResponse(nodes))
}

func SetDataSourceNodeRole(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	if params["role"] != "master" {
		return errors.New("you can only set a node to master")
	}

	node := GetDataSourceNodeInfo(ctx)

	err = node.SetMaster()

	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func CreateDataSourceNode(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	dataSourceNode := &models.DataSourceNode{
		Host: params["host"],
		Role: "master",
	}

	dataSource := GetDataSourceInfo(ctx)

	err = dataSource.CreateNode(dataSourceNode)

	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewDataSourceNodeResponse(dataSourceNode))
}

func UpdateDataSourceNode(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	updates := &models.DataSourceNode{
		Role:       params["role"],
		MasterHost: params["masterHost"],
	}

	node := GetDataSourceNodeInfo(ctx)

	err = node.Update(updates)

	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}

func PollDataSourceNodeCommands(ctx echo.Context) error {
	node := GetDataSourceNodeInfo(ctx)

	command, err := node.WaitForCommand()

	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, command)
}
