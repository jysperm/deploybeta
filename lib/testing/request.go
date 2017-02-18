package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo"
)

func RequestJSON(handler echo.HandlerFunc, method string, url string, body interface{}) (*httptest.ResponseRecorder, error) {
	app := echo.New()

	json, err := json.Marshal(body)

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(json))

	if err != nil {
		panic(err)
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	res := httptest.NewRecorder()
	ctx := app.NewContext(req, res)
	err = handler(ctx)

	return res, err
}
