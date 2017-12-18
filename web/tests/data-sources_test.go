package tests

import (
	"strings"
	"testing"

	. "github.com/jysperm/deploying/lib/testing"
	"github.com/jysperm/deploying/lib/utils"
)

func TestCreateDataSource(t *testing.T) {
	res, _, errs := Request("POST", "/data-sources").
		Set("Authorization", globalSession.Token).
		SendStruct(map[string]string{
			"name": strings.ToLower(utils.RandomString(10)),
			"type": "redis",
		}).EndBytes()

	if res.StatusCode != 201 || len(errs) != 0 {
		t.Error(errs)
	}
}
