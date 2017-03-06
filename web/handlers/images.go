package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/docker/docker/api/types"

	"github.com/jysperm/deploying/lib/services/builder"
	. "github.com/jysperm/deploying/web/handlers/helpers"
)

func generateVersion() string {
	now := time.Now()
	return fmt.Sprintf("%d%d%d-%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func CreateImage(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	version := generateVersion()
	buildOpts := types.ImageBuildOptions{
		Tags: []string{version},
	}
	shasum, err := builder.BuildImage(buildOpts, params["gitRepository"])
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewImageResponse(shasum, version))
}
