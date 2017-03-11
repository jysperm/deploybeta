package tests

import (
	"fmt"
	"testing"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	appModel "github.com/jysperm/deploying/lib/models/app"
	sessionModel "github.com/jysperm/deploying/lib/models/session"
	. "github.com/jysperm/deploying/lib/testing"
	"github.com/jysperm/deploying/web/handlers/helpers"
)

func TestCreateImage(t *testing.T) {
	account, _ := SeedAccount()
	session := SeedSession(&account)
	app := SeedApp("https://github.com/mason96112569/docker-test.git")

	var imageResponse helpers.ImageResponse

	imageKey := fmt.Sprintf("/apps/%s/images", app.Name)
	res, _, errs := Request("POST", imageKey).
		Set("Authorization", session.Token).
		SendStruct(app).
		EndStruct(&imageResponse)

	t.Log(imageResponse)

	if len(errs) != 0 {
		t.Error(errs)
	}

	if res.StatusCode != 201 {
		t.Error("res.StatusCode", res.StatusCode)
	}

	accountModel.DeleteByName(session.Username)
	sessionModel.DeleteByToken(session.Token)
	appModel.DeleteByName(app.Name)
}
