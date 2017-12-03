package handlers

import (
	"testing"

	"github.com/jysperm/deploying/lib/models"
	. "github.com/jysperm/deploying/lib/testing"
	"github.com/labstack/echo"
)

func TestCreateSession(t *testing.T) {
	account, password := SeedAccount()

	res, body, err := RequestJSON(CreateSession, echo.POST, "/sessions", map[string]string{
		"username": account.Username,
		"password": password,
	})

	if err != nil {
		t.Error(err)
	}

	if res.Code != 201 {
		t.Errorf("res.Code %v", res.Code)
	}

	models.DeleteAccountByName(account.Username)
	models.DeleteSessionByToken(body["token"])
}
