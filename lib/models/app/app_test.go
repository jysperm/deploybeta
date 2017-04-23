package app

import (
	"strings"
	"testing"

	"github.com/jysperm/deploying/lib/utils"
)

func TestUpdateVersion(t *testing.T) {
	app := Application{
		Name: strings.ToLower(utils.RandomString(10)),
	}

	err := CreateApp(&app)

	if err != nil {
		panic(err)
	}

	err = app.UpdateVersion("20170411-212400")

	if err != nil {
		panic(err)
	}

	if app.Version != "20170411-212400" {
		t.Errorf("app.Version %s", app.Version)
	}

	latestApp, err := FindByName(app.Name)

	if err != nil {
		panic(err)
	}

	if latestApp.Version != "20170411-212400" {
		t.Errorf("latestApp.Version %s", latestApp.Version)
	}

	DeleteByName(app.Name)
}
