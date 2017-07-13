package node

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jysperm/deploying/lib/utils"
)

func init() {
	GOPATH := os.Getenv("GOPATH")
	root := filepath.Join(GOPATH, "src", "github.com", "jysperm", "deploying")
	if err := os.Setenv("WORKDIR", root); err != nil {
		panic(err)
	}
}

func TestYarn(t *testing.T) {
	root, err := utils.Clone("https://github.com/jysperm/deploying-samples.git", "yarn")
	if err != nil {
		t.Error(err)
	}
	if err := GenerateDockerfile("", root); err != nil {
		t.Error(err)
	}
	dockerfilePath := filepath.Join(root, "Dockerfile")
	content, err := ioutil.ReadFile(dockerfilePath)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(content))
}

func TestNpm(t *testing.T) {
	root, err := utils.Clone("https://github.com/jysperm/deploying-samples.git", "npm")
	if err != nil {
		t.Error(err)
	}

	dockerfilePath := filepath.Join(root, "Dockerfile")
	content, err := ioutil.ReadFile(dockerfilePath)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(content))
}
