package tests

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
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

func TestMain(m *testing.M) {
	flag.Parse()
	globalAccount, _ = SeedAccount()
	globalSession = SeedSession(&globalAccount)
	globalApp = SeedApp("https://github.com/jysperm/deploying-samples.git", globalAccount.Username)

	exitVal := m.Run()

	sessionModel.DeleteByToken(globalSession.Token)
	accountModel.DeleteByName(globalAccount.Username)
	swarm.RemoveService(globalApp)
	os.Exit(exitVal)
}

func TestCreateVersion(t *testing.T) {
	var version versionModel.Version
	requestPath := fmt.Sprintf("/apps/%s/versions", globalApp.Name)
	res, body, errs := Request("POST", requestPath).
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"gitTag": "master",
		}).EndBytes()
	if res.StatusCode != 201 || len(errs) != 0 {
		t.Error(errs)
	}
	if err := json.Unmarshal(body, &version); err != nil {
		t.Error(err)
	}
	t.Log("Created version: ", version)
}

func TestDeployVersion(t *testing.T) {
	var version versionModel.Version

	requestPath := fmt.Sprintf("/apps/%s/versions", globalApp.Name)
	res, body, errs := Request("POST", requestPath).
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"gitTag": "master",
		}).EndBytes()
	if res.StatusCode != 201 || len(errs) != 0 {
		t.Error(errs)
	}

	if err := json.Unmarshal(body, &version); err != nil {
		t.Error(err)
	}

	deployPath := fmt.Sprintf("/apps/%s/version", globalApp.Name)
	res, _, errs = Request("PUT", deployPath).
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"tag": version.Tag,
		}).EndBytes()
	if res.StatusCode != 200 || len(errs) != 0 {
		t.Error(errs)
	}

	t.Log("Deployed version: ", version)
}

func TestCreateAndDeploy(t *testing.T) {
	var version versionModel.Version
	requestPath := fmt.Sprintf("/apps/%s/version", globalApp.Name)
	res, body, errs := Request("POST", requestPath).
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"gitTag": "master",
		}).EndBytes()
	if res.StatusCode != 201 || len(errs) != 0 {
		t.Error(errs)
	}

	if err := json.Unmarshal(body, &version); err != nil {
		t.Error(err)
	}

	t.Log("Created and Deployed version: ", version)
}
