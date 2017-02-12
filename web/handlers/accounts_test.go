package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jysperm/deploying/lib/utils"
	"github.com/labstack/echo"
)

var app = echo.New()

func TestRegisterAccountSuccess(t *testing.T) {
	json, err := json.Marshal(map[string]string{
		"username": utils.RandomString(10),
		"email":    utils.RandomString(10) + "@gmail.com",
		"password": utils.RandomString(10),
	})

	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(echo.POST, "/accounts", bytes.NewBuffer(json))

	if err != nil {
		t.Error(err)
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	res := httptest.NewRecorder()
	ctx := app.NewContext(req, res)

	err = RegisterAccount(ctx)

	if err != nil {
		t.Error(err)
	}

	if res.Code != 201 {
		t.Errorf("res.Code %v", res.Code)
	}
}
