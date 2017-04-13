package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	appModel "github.com/jysperm/deploying/lib/models/app"
	sessionModel "github.com/jysperm/deploying/lib/models/session"
	versionModel "github.com/jysperm/deploying/lib/models/version"
	"github.com/jysperm/deploying/lib/swarm"
	. "github.com/jysperm/deploying/lib/testing"
)

var globalAccount accountModel.Account
var globalSession sessionModel.Session
var globalApp appModel.Application
var globalVersion versionModel.Version

func TestCreateVersion(t *testing.T) {
	globalAccount, _ = SeedAccount()
	globalSession = SeedSession(&globalAccount)
	globalApp = SeedApp("https://github.com/jysperm/deploying-samples.git")

	requestPath := fmt.Sprintf("/apps/%s/versions", globalApp.Name)
	res, body, errs := Request("POST", requestPath).
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"gitTag": "master",
		}).EndBytes()
	if res.StatusCode != 201 || len(errs) != 0 {
		t.Error(errs)
	}
	if err := json.Unmarshal(body, &globalVersion); err != nil {
		t.Error(err)
	}
}

func TestDeployVersion(t *testing.T) {
	deployPath := fmt.Sprintf("/apps/%s/version", globalApp.Name)
	res, _, errs := Request("PUT", deployPath).
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"tag": globalVersion.Tag,
		}).EndBytes()
	if res.StatusCode != 201 || len(errs) != 0 {
		t.Error(errs)
	}

	accountModel.DeleteByName(globalSession.Username)
	sessionModel.DeleteByToken(globalSession.Token)
	versionModel.DeleteVersion(globalApp, globalVersion.Tag)
	appModel.DeleteByName(globalApp.Name)
	swarm.RemoveService(globalApp)
}

func TestCreateAndDeploy(t *testing.T) {
	account, _ := SeedAccount()
	session := SeedSession(&account)
	app := SeedApp("https://github.com/jysperm/deploying-samples.git")

	var appVersion versionModel.Version
	requestPath := fmt.Sprintf("/apps/%s/version", app.Name)
	res, body, errs := Request("POST", requestPath).
		Set("Authorization", session.Token).
		SendStruct(map[string]string{
			"gitTag": "master",
		}).EndBytes()
	if res.StatusCode != 201 || len(errs) != 0 {
		t.Error(errs)
	}
	if err := json.Unmarshal(body, &appVersion); err != nil {
		t.Error(err)
	}

	accountModel.DeleteByName(session.Username)
	sessionModel.DeleteByToken(session.Token)
	versionModel.DeleteVersion(app, appVersion.Tag)
	appModel.DeleteByName(app.Name)
	swarm.RemoveService(app)
}
