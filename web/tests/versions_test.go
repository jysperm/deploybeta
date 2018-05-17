package tests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/jysperm/deploybeta/lib/models"
	. "github.com/jysperm/deploybeta/lib/testing"
)

func TestCreateVersion(t *testing.T) {
	var version models.Version
	requestPath := fmt.Sprintf("/apps/%s/versions", globalApp.Name)
	res, body, errs := Request("POST", requestPath).
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"gitTag": "dep",
		}).EndBytes()
	if res.StatusCode != 201 || len(errs) != 0 {
		t.Error(errs)
	}
	if err := json.Unmarshal(body, &version); err != nil {
		t.Error(err)
	}

	progressPath := fmt.Sprintf("http://127.0.0.1:7000/apps/%s/versions/%s/progress", globalApp.Name, version.Tag)
	client := &http.Client{}
	req, err := http.NewRequest("GET", progressPath, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", globalSession.Token)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	for {
		line, _ := reader.ReadBytes('\n')
		strLine := string(line)
		if strings.Contains(strLine, "Deploybeta: Building finished.") {
			break
		}
		fmt.Print(string(line))
	}

	v, err := models.FindVersionByTag(&globalApp, version.Tag)
	if err != nil {
		t.Error(err)
	}
	t.Log("Created version: ", v)
}

func TestDeployVersion(t *testing.T) {
	var version models.Version

	requestPath := fmt.Sprintf("/apps/%s/versions", globalApp.Name)
	res, body, errs := Request("POST", requestPath).
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"gitTag": "yarn",
		}).EndBytes()
	if res.StatusCode != 201 || len(errs) != 0 {
		t.Error(errs)
	}

	if err := json.Unmarshal(body, &version); err != nil {
		t.Error(err)
	}

	progressPath := fmt.Sprintf("http://127.0.0.1:7000/apps/%s/versions/%s/progress", globalApp.Name, version.Tag)
	client := &http.Client{}
	req, err := http.NewRequest("GET", progressPath, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", globalSession.Token)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	for {
		line, _ := reader.ReadBytes('\n')
		strLine := string(line)
		if strings.Contains(strLine, "Deploybeta: Building finished.") {
			break
		}
		fmt.Print(string(line))
	}

	deployPath := fmt.Sprintf("/apps/%s/version", globalApp.Name)
	res, _, errs = Request("PUT", deployPath).
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"tag": version.Tag,
		}).EndBytes()
	if res.StatusCode != 200 || len(errs) != 0 {
		t.Fatal(errs)
	}

}

func TestPushProgress(t *testing.T) {
	var version models.Version
	requestPath := fmt.Sprintf("/apps/%s/versions", globalApp.Name)
	res, body, errs := Request("POST", requestPath).
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"gitTag": "npm",
		}).EndBytes()
	if res.StatusCode != 201 || len(errs) != 0 {
		t.Fatal(errs)
	}
	if err := json.Unmarshal(body, &version); err != nil {
		t.Error(err)
	}

	progressPath := fmt.Sprintf("http://127.0.0.1:7000/apps/%s/versions/%s/progress", globalApp.Name, version.Tag)
	client := &http.Client{}
	req, err := http.NewRequest("GET", progressPath, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", globalSession.Token)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	for {
		line, _ := reader.ReadBytes('\n')
		strLine := string(line)
		if strings.Contains(strLine, "Deploybeta: Building finished.") {
			break
		}
		fmt.Print(string(line))
	}

}
