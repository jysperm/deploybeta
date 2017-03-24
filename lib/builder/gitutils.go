package builder

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/docker/docker/pkg/urlutil"
)

func Clone(remoteURL string, param string) (string, error) {
	if !urlutil.IsGitURL(remoteURL) {
		return "", errors.New("Not a valid git URL")
	}

	root, err := ioutil.TempDir("", "deploying-build")
	if err != nil {
		return "", err
	}

	if output, err := git("clone", remoteURL, "--branch", param, root); err != nil {
		return "", fmt.Errorf("Error trying to use git: %s (%s) ", err, output)
	}

	return root, nil
}

func git(args ...string) ([]byte, error) {
	return exec.Command("git", args...).CombinedOutput()
}
