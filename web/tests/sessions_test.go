package tests

import (
	"encoding/json"
	"testing"

	"github.com/jysperm/deploybeta/lib/models"
	. "github.com/jysperm/deploybeta/lib/testing"
)

func TestCreateSession(t *testing.T) {
	account, password := SeedAccount()

	res, body, errs := Request("POST", "/sessions").
		SendStruct(map[string]string{
			"username": account.Username,
			"password": password,
		}).EndBytes()

	if res.StatusCode != 201 || len(errs) != 0 {
		t.Error(res.StatusCode, errs)
	}

	session := &models.Session{}

	if err := json.Unmarshal(body, session); err != nil {
		t.Error(err)
	}

	if err := session.Destroy(); err != nil {
		t.Error(err)
	}

	if err := account.Destroy(); err != nil {
		t.Error(err)
	}
}
