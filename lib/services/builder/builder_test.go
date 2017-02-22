package builder

import (
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestBuildImage(t *testing.T) {

	opts := types.ImageBuildOptions{
		Tags:       []string{"docker-test"},
		Dockerfile: "Dockerfile",
	}

	shasum, err := BuildImage(opts, "https://github.com/mason96112569/docker-test.git")
	if err != nil {
		t.Error(err)
	}
	t.Log(shasum)
}

func TestBuildFailure(t *testing.T) {
	opts := types.ImageBuildOptions{
		Tags:       []string{"failure"},
		Dockerfile: "Dockerfile",
	}

	_, err := BuildImage(opts, "https://github.com/mason96112569/docker-test-failure.git")
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Error(errors.New("It should have a error"))
	}

}
