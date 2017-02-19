package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	sessionModel "github.com/jysperm/deploying/lib/models/session"
	. "github.com/jysperm/deploying/lib/testing"
	"github.com/jysperm/deploying/lib/utils"
	"github.com/jysperm/deploying/web/handlers/helpers"
)

func TestRegisterAccount(t *testing.T) {
	username := utils.RandomString(10)

	res, _, err := RequestJSON(RegisterAccount, echo.POST, "/accounts", map[string]string{
		"username": username,
		"email":    utils.RandomString(10) + "@gmail.com",
		"password": utils.RandomString(10),
	})

	if err != nil {
		t.Error(err)
	}

	if res.Code != 201 {
		t.Errorf("res.Code %v", res.Code)
	}

	res, _, err = RequestJSON(RegisterAccount, echo.POST, "/accounts", map[string]string{
		"username": username,
		"email":    utils.RandomString(10) + "@gmail.com",
		"password": utils.RandomString(10),
	})

	if err == nil || err.(*echo.HTTPError).Code != 409 {
		t.Errorf("err.Code %v", err.(*echo.HTTPError).Code)
	}

	accountModel.DeleteByName(username)
}

func TestCurrentAccount(t *testing.T) {
	app := echo.New()

	account, _ := SeedAccount()
	session := SeedSession(&account)

	req, err := http.NewRequest(echo.GET, "/session/account", nil)

	if err != nil {
		t.Error(err)
	}

	req.Header.Set(echo.HeaderAuthorization, session.Token)

	res := httptest.NewRecorder()
	ctx := app.NewContext(req, res)

	err = helpers.AuthenticateMiddleware(CurrentAccount)(ctx)

	if err != nil {
		t.Error(err)
	}

	if res.Code != 200 {
		t.Errorf("res.Code %v", res.Code)
	}

	response := map[string]string{}
	err = json.Unmarshal(res.Body.Bytes(), &response)

	if err != nil {
		t.Error(err)
	}

	if response["username"] != account.Username {
		t.Errorf("response.username %v", response["username"])
	}

	accountModel.DeleteByName(account.Username)
	sessionModel.DeleteByToken(session.Token)
}
