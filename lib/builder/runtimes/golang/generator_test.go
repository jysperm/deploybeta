package golang

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jysperm/deploying/lib/builder"
)

func init() {
	IsTesting = true
}

func TestGenerateDockerfile(t *testing.T) {
	depRoot, err := builder.Clone("https://github.com/jysperm/deploying-samples.git", "dep")
	if err != nil {
		t.Error(err)
	}
	if err := GenerateDockerfile(depRoot, Go181, "github.com/jysperm/deploying-samples", "deploying-samples"); err != nil {
		t.Error(err)
	}
	dockerfilePath := filepath.Join(depRoot, "Dockerfile")
	depDockerfile, err := ioutil.ReadFile(dockerfilePath)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(depDockerfile))

	glideRoot, err := builder.Clone("https://github.com/jysperm/deploying-samples.git", "glide")
	if err != nil {
		t.Error(err)
	}
	if err := GenerateDockerfile(glideRoot, Go181, "github.com/jysperm/deploying-samples", "deploying-samples"); err != nil {
		t.Error(err)
	}
	dockerfilePath = filepath.Join(glideRoot, "Dockerfile")
	glideDockerfile, err := ioutil.ReadFile(dockerfilePath)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(glideDockerfile))
}
