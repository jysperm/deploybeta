package tests

import (
	"fmt"
	"testing"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	appModel "github.com/jysperm/deploying/lib/models/app"
	sessionModel "github.com/jysperm/deploying/lib/models/session"
	versionModel "github.com/jysperm/deploying/lib/models/version"
	. "github.com/jysperm/deploying/lib/testing"
)

func TestCreateImage(t *testing.T) {
	account, _ := SeedAccount()
	session := SeedSession(&account)
	app := SeedApp("https://github.com/mason96112569/docker-test.git")

	var imageResponse versionModel.Version

	imageKey := fmt.Sprintf("/apps/%s/images", app.Name)
	res, _, errs := Request("POST", imageKey).
		Set("Authorization", session.Token).
		SendStruct(map[string]string{
			"name": app.Name,
		}).
		EndStruct(&imageResponse)

	if len(errs) != 0 {
		t.Error(errs)
	}

	if res.StatusCode != 200 {
		t.Error("res.StatusCode", res.StatusCode)
	}

	t.Log(imageResponse)
	accountModel.DeleteByName(session.Username)
	sessionModel.DeleteByToken(session.Token)
	versionModel.DeleteVersion(app, imageResponse.Tag)
	appModel.DeleteByName(app.Name)
}
