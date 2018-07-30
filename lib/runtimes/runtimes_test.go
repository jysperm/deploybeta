package runtimes

import (
	"testing"

	"github.com/jysperm/deploybeta/lib/utils"
)

const sampleGitUrl = "https://github.com/jysperm/deploybeta-samples.git"

func TestNodejs(t *testing.T) {
	dockerlize(t, "nodejs-npm-express")
	dockerlize(t, "nodejs-yarn-express")
}

func TestGolang(t *testing.T) {
	dockerlize(t, "golang-dep-echo")
	dockerlize(t, "golang-glide-echo")
}

func dockerlize(t *testing.T, branch string) {
	root, err := utils.Clone(sampleGitUrl, branch)

	if err != nil {
		t.Fatal(err)
	}

	runtime, err := DecideRuntime(NewBuildContext(root, sampleGitUrl))

	if err != nil {
		t.Fatal(err)
	}

	dockerfile, err := runtime.Dockerfile()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(branch)
	t.Log(dockerfile.String())
}
