package tests

import (
	"fmt"
	"strings"
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
	appName := strings.ToLower(utils.RandomString(10))

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

func TestUpdateApp(t *testing.T) {
	account, _ := SeedAccount()
	session := SeedSession(&account)
	app := SeedApp("https://github.com/jysperm/deploying-samples.git", account.Username)

	update := appModel.Application{
		Name:          app.Name,
		Instances:     3,
		GitRepository: "https://github.com/jysperm/deploying-samples.git",
	}

	newApp := appModel.Application{}

	updateURL := fmt.Sprintf("/apps/%s", app.Name)
	res, _, errs := Request("PATCH", updateURL).
		Set("Authorization", session.Token).
		SendStruct(update).
		EndStruct(&newApp)

	if len(errs) != 0 {
		t.Error(errs)
	}

	if res.StatusCode != 201 {
		t.Error("Updateing failed.")
	}

	t.Log(newApp)
}

func TestDeleteApp(t *testing.T) {
	account, _ := SeedAccount()
	session := SeedSession(&account)
	app := SeedApp("https://github.com/jysperm/deploying-samples.git", account.Username)

	deleteURL := fmt.Sprintf("/apps/%s", app.Name)

	res, _, errs := Request("DELETE", deleteURL).
		Set("Authorization", session.Token).
		SendStruct(map[string]string{
			"name": app.Name,
		}).End()

	if len(errs) != 0 {
		t.Error(errs)
	}

	if res.StatusCode != 200 {
		t.Error("Deleting app failed")
	}
}
