package tests

import (
	"testing"

	"github.com/jysperm/deploying/config"
	accountModel "github.com/jysperm/deploying/lib/models/account"
	appModel "github.com/jysperm/deploying/lib/models/app"
	sessionModel "github.com/jysperm/deploying/lib/models/session"
	. "github.com/jysperm/deploying/lib/testing"
	"github.com/jysperm/deploying/lib/utils"
	"github.com/jysperm/deploying/web"
)

func init() {
	go web.CreateWebServer().Start(config.Port)
}

func TestCreateApp(t *testing.T) {
	account, _ := SeedAccount()
	session := SeedSession(&account)
	appName := utils.RandomString(10)

	resBody := appModel.Application{}

	res, _, errs := Request("POST", "/apps").
		Set("Authorization", session.Token).
		SendStruct(map[string]string{
			"name": appName,
		}).EndStruct(&resBody)

	t.Log("Created app", resBody)

	if len(errs) != 0 {
		t.Error(errs)
	}

	if res.StatusCode != 201 {
		t.Errorf("res.StatusCode %v", res.StatusCode)
	}

	res, _, errs = Request("POST", "/apps").
		Set("Authorization", session.Token).
		SendStruct(map[string]string{
			"name": appName,
		}).End()

	if len(errs) != 0 {
		t.Error(errs)
	}

	if res.StatusCode != 409 {
		t.Errorf("res.StatusCode %v", res.StatusCode)
	}

	apps := []appModel.Application{}

	res, _, errs = Request("GET", "/apps").
		Set("Authorization", session.Token).
		EndStruct(&apps)

	if len(errs) != 0 {
		t.Error(errs)
	}

	if apps[0].Name != appName {
		t.Errorf("appName %v", apps[0].Name)
	}

	accountModel.DeleteByName(session.Username)
	sessionModel.DeleteByToken(session.Token)
	appModel.DeleteByName(appName)
}
