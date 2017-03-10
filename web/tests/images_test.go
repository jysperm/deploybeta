package tests

import (
	"fmt"
	"testing"

	"github.com/jysperm/deploying/config"
	accountModel "github.com/jysperm/deploying/lib/models/account"
	appModel "github.com/jysperm/deploying/lib/models/app"
	sessionModel "github.com/jysperm/deploying/lib/models/session"
	. "github.com/jysperm/deploying/lib/testing"
	"github.com/jysperm/deploying/web"
)

func init() {
	go web.CreateWebServer().Start(config.Port)
}

func TestCreateImage(t *testing.T) {
	account, _ := SeedAccount()
	session := SeedSession(&account)
	app := SeedApp("https://github.com/jysperm/deploying-samples.git")

	resBody := appModel.Application{}

	res, _, errs := Request("POST", "/apps").
		Set("Authorization", session.Token).
		SendStruct(map[string]string{
			"name": app.Name,
		}).EndStruct(&resBody)

	t.Log("Created app", resBody)

	if len(errs) != 0 {
		t.Error(errs)
	}

	if res.StatusCode != 201 {
		t.Errorf("res.StatusCode %v", res.StatusCode)

	}

	imageKey := fmt.Sprintf("/apps/%s/images", app.Name)
	res, _, errs = Request("POST", imageKey).
		Set("Authorization", session.Token).
		EndStruct(app)

	if len(errs) != 0 {
		t.Error(errs)
	}

	accountModel.DeleteByName(session.Username)
	sessionModel.DeleteByToken(session.Token)
	appModel.DeleteByName(app.Name)
}
