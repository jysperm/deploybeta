package builder

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os/exec"

	"github.com/docker/docker/pkg/urlutil"
)

func Clone(remoteURL string, param string) (string, error) {
	if !urlutil.IsGitURL(remoteURL) {
		return "", errors.New("Not a valid git URL")
	}
	if !urlutil.IsTransportURL(remoteURL) {
		remoteURL = "https://" + remoteURL
	}

	root, err := ioutil.TempDir("", "deploying-build")
	if err != nil {
		return "", err
	}

	u, err := url.Parse(remoteURL)
	if err != nil {
		return "", err
	}

	if u.Fragment != "" {
		u.Fragment = ""
	}

	args := fmt.Sprintf(`clone --branch %s %s %s`, param, remoteURL, root)

	if output, err := git(args); err != nil {
		return "", fmt.Errorf("Error trying to use git: %s (%s) ", err, output)
	}

	return root, nil
}

func git(args ...string) ([]byte, error) {
	return exec.Command("git", args...).CombinedOutput()
}
