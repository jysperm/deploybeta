package builder

import (
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestBuildDep(t *testing.T) {
	opts := types.ImageBuildOptions{
		Tags: []string{"dep:latest"},
	}

	shasum, err := BuildImage(opts, "https://github.com/jysperm/deploying-samples.git", "dep")
	if err != nil {
		t.Error(err)
	}
	t.Log(shasum)
}

func TestBuildGlide(t *testing.T) {
	opts := types.ImageBuildOptions{
		Tags: []string{"glide:latest"},
	}

	shasum, err := BuildImage(opts, "https://github.com/jysperm/deploying-samples.git", "glide")
	if err != nil {
		t.Error(err)
	}
	t.Log(shasum)
}

func TestBuildNpm(t *testing.T) {
	opts := types.ImageBuildOptions{
		Tags: []string{"npm:latest"},
	}

	shasum, err := BuildImage(opts, "https://github.com/jysperm/deploying-samples.git", "npm")
	if err != nil {
		t.Error(err)
	}
	t.Log(shasum)
}

func TestBuildYarn(t *testing.T) {
	opts := types.ImageBuildOptions{
		Tags: []string{"yarn:latest"},
	}

	shasum, err := BuildImage(opts, "https://github.com/jysperm/deploying-samples.git", "yarn")
	if err != nil {
		t.Error(err)
	}
	t.Log(shasum)
}
func TestBuildUnknownProject(t *testing.T) {
	opts := types.ImageBuildOptions{
		Tags: []string{"failure:latest"},
	}

	_, err := BuildImage(opts, "https://github.com/jysperm/deploying-samples.git", "failure")
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Error(errors.New("It should have a error"))
	}

}
