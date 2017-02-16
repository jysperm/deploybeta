package builder

import (
	"testing"

	"github.com/docker/docker/api/types"
)

func TestBuildImage(t *testing.T) {

	opts := types.ImageBuildOptions{
		Tags:       []string{"docker-test"},
		Dockerfile: "Dockerfile",
	}

	id, err := BuildImage(opts, "https://github.com/mason96112569/docker-test.git")
	if err != nil {
		t.Error(err)
	}
	t.Log("Image ID is " + id)
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
		t.Error("It should have an error")
	}

}
