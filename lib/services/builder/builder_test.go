package builder

import (
	"testing"

	"github.com/docker/docker/api/types"
)

func TestBuildImage(t *testing.T) {
	opts := types.ImageBuildOptions{
		Tags:           []string{"docker-test", "ubuntu"},
		SuppressOutput: true,
		NoCache:        false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Isolation:      "",
	}

	_, err := BuildImage(opts, "https://github.com/mason96112569/docker-test.git")
	if err != nil {
		t.Error(err)
	}

}
