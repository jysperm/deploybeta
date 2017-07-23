package golang

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
func TestGenerateDockerfile(t *testing.T) {
	depRoot, err := utils.Clone("https://github.com/jysperm/deploying-samples.git", "dep")
	if err != nil {
		t.Error(err)
	}
	if err := GenerateDockerfile(depRoot, "https://github.com/jysperm/deploying-samples.git"); err != nil {
		t.Error(err)
	}
	dockerfilePath := filepath.Join(depRoot, "Dockerfile")
	depDockerfile, err := ioutil.ReadFile(dockerfilePath)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(depDockerfile))

	glideRoot, err := utils.Clone("git@github.com:jysperm/deploying-samples.git", "glide")
	if err != nil {
		t.Error(err)
	}
	if err := GenerateDockerfile(glideRoot, "git@github.com:jysperm/deploying-samples.git"); err != nil {
		t.Error(err)
	}
	dockerfilePath = filepath.Join(glideRoot, "Dockerfile")
	glideDockerfile, err := ioutil.ReadFile(dockerfilePath)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(glideDockerfile))
}