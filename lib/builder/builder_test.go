package builder

import (
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestBuildImage(t *testing.T) {

	opts := types.ImageBuildOptions{
		Tags: []string{"dep-test:latest"},
	}

	shasum, err := BuildImage(opts, "https://github.com/jysperm/deploying-samples.git", "dep")
	if err != nil {
		t.Error(err)
	}
	t.Log(shasum)
}

func TestBuildFailure(t *testing.T) {
	opts := types.ImageBuildOptions{
		Tags: []string{"failure"},
	}

	_, err := BuildImage(opts, "https://github.com/jysperm/deploying-samples.git", "failure")
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Error(errors.New("It should have a error"))
	}

}
