package tests

import (
	"strings"
	"testing"

	. "github.com/jysperm/deploying/lib/testing"
	"github.com/jysperm/deploying/lib/utils"
)

var dataSourceName string

func TestCreateDataSource(t *testing.T) {
	dataSourceName = strings.ToLower(utils.RandomString(10))

	res, _, errs := Request("POST", "/data-sources").
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"name": dataSourceName,
			"type": "redis",
		}).EndBytes()

	if res.StatusCode != 201 || len(errs) != 0 {
		t.Error(errs)
	}
}

func TestListDataSources(t *testing.T) {
	res, _, errs := Request("GET", "/data-sources").
		Set("Authorization", globalSession.Token).
		EndBytes()

	if res.StatusCode != 200 || len(errs) != 0 {
		t.Error(errs)
	}
}

func TestCreateDataSourceNode(t *testing.T) {
	res, _, errs := Request("POST", "/data-sources/"+dataSourceName+"/agents").
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"host": "127.0.0.1",
		}).EndBytes()

	if res.StatusCode != 201 || len(errs) != 0 {
		t.Error(errs)
	}
}
