package handlers

import (
	"testing"

	"github.com/labstack/echo"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	. "github.com/jysperm/deploying/lib/testing"
	"github.com/jysperm/deploying/lib/utils"
)

var app = echo.New()

func TestRegisterAccount(t *testing.T) {
	username := utils.RandomString(10)

	res, err := RequestJSON(RegisterAccount, echo.POST, "/accounts", map[string]string{
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

	res, err = RequestJSON(RegisterAccount, echo.POST, "/accounts", map[string]string{
		"username": username,
		"email":    utils.RandomString(10) + "@gmail.com",
		"password": utils.RandomString(10),
	})

	if err == nil || err.(*echo.HTTPError).Code != 409 {
		t.Errorf("err.Code %v", err.(*echo.HTTPError).Code)
	}

	accountModel.DeleteByName(username)
}
