package tests

import (
	"encoding/json"
	"testing"

	"github.com/jysperm/deploybeta/lib/models"
	. "github.com/jysperm/deploybeta/lib/testing"
	"github.com/jysperm/deploybeta/lib/utils"
	"github.com/jysperm/deploybeta/web/handlers/helpers"
)

func TestRegisterAccount(t *testing.T) {
	username := utils.RandomString(10)

	res, _, errs := Request("POST", "/accounts").
		SendStruct(map[string]string{
			"username": username,
			"email":    utils.RandomString(10) + "@gmail.com",
			"password": utils.RandomString(10),
		}).EndBytes()

	if res.StatusCode != 201 || len(errs) != 0 {
		t.Error(res.StatusCode, errs)
	}

	res, _, errs = Request("POST", "/accounts").
		SendStruct(map[string]string{
			"username": username,
			"email":    utils.RandomString(10) + "@gmail.com",
			"password": utils.RandomString(10),
		}).EndBytes()

	if res.StatusCode != 409 || len(errs) != 0 {
		t.Error(res.StatusCode, errs)
	}

	if err := (&models.Account{Username: username}).Destroy(); err != nil {
		t.Error(err)
	}
}

func TestCurrentAccount(t *testing.T) {
	res, body, errs := Request("GET", "/session/account").
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"name": dataSourceName,
			"type": "redis",
		}).EndBytes()

	if res.StatusCode != 200 || len(errs) != 0 {
		t.Error(res.StatusCode, errs)
	}

	response := &helpers.AccountResponse{}

	if err := json.Unmarshal(body, response); err != nil {
		t.Error(err)
	}

	if response.Username != globalAccount.Username {
		t.Error(response)
	}
}
