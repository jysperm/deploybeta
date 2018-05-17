package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
)

var validCommitHash = regexp.MustCompile("^[0-9a-f]{7,40}$")
var validGitURL = regexp.MustCompile(`(^(https|http)://[^\s]*.git$)|(^\w+@\w+\.\w+:[^\s]*.git$)`)

func Clone(remoteURL string, gitTag string) (string, error) {
	if !validGitURL.MatchString(remoteURL) {
		return "", errors.New("Not a valid git URL")
	}

	root, err := ioutil.TempDir("", "deploybeta-build")
	if err != nil {
		return "", err
	}

	if gitTag == "" {
		gitTag = "master"
	}

	if !validCommitHash.MatchString(gitTag) {
		if output, err := git("clone", remoteURL, "--branch", gitTag, root); err != nil {
			return "", fmt.Errorf("Error trying to use git: %s (%s) ", err, output)
		}
		return root, nil
	}

	if output, err := git("clone", remoteURL, root); err != nil {
		return "", fmt.Errorf("Error trying to use git: %s (%s) ", err, output)
	}
	if output, err := git("checkout", gitTag, "-C", root); err != nil {
		return "", fmt.Errorf("Error trying to use git: %s (%s) ", err, output)
	}

	return root, nil
}

func git(args ...string) ([]byte, error) {
	return exec.Command("git", args...).CombinedOutput()
}
