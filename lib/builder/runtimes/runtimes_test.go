package runtimes

import (
	"testing"

	"github.com/jysperm/deploying/lib/utils"
)

func TestDockerlizeDep(t *testing.T) {
	root, err := utils.Clone("https://github.com/jysperm/deploying-samples.git", "dep")
	err = Dockerlize(root, "https://github.com/jysperm/deploying-samples.git")
	if err != nil {
		t.Error(err)
	}
}

func TestDockerlizeGlide(t *testing.T) {
	root, err := utils.Clone("https://github.com/jysperm/deploying-samples.git", "glide")
	err = Dockerlize(root, "https://github.com/jysperm/deploying-samples.git")
	if err != nil {
		t.Error(err)
	}
}
func TestDockerlizeNpm(t *testing.T) {
	root, err := utils.Clone("https://github.com/jysperm/deploying-samples.git", "npm")
	err = Dockerlize(root, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestDockerlizeYarn(t *testing.T) {
	root, err := utils.Clone("https://github.com/jysperm/deploying-samples.git", "yarn")
	err = Dockerlize(root, nil)
	if err != nil {
		t.Error(err)
	}
}
