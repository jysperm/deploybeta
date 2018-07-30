package runtimes

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jysperm/deploybeta/config"
)

type BuildContext struct {
	root   string
	gitUrl string

	ProxyCommand string
	AptMirror    string
}

func NewBuildContext(root string, gitUrl string) *BuildContext {
	return &BuildContext{
		root:         root,
		gitUrl:       gitUrl,
		ProxyCommand: getProxyCommand(),
		AptMirror:    config.AptMirror,
	}
}

func (context *BuildContext) fileExists(name string) (bool, error) {
	filename := filepath.Join(context.root, name)

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func (context *BuildContext) readFile(name string) ([]byte, error) {
	filename := filepath.Join(context.root, name)

	fileContent, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

func getProxyCommand() string {
	return "http_proxy=" + config.HttpProxy + " https_proxy=" + config.HttpProxy
}
