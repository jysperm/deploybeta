package runtimes

import (
	"testing"

	"github.com/jysperm/deploybeta/lib/utils"
)

func TestDockerlizeDep(t *testing.T) {
	root, err := utils.Clone("https://github.com/jysperm/deploying-samples.git", "dep")
	buf, err := Dockerlize(root, "https://github.com/jysperm/deploying-samples.git")
	if err != nil {
		t.Error(err)
	}
	t.Log(buf.String())
}

func TestDockerlizeGlide(t *testing.T) {
	root, err := utils.Clone("https://github.com/jysperm/deploying-samples.git", "glide")
	buf, err := Dockerlize(root, "https://github.com/jysperm/deploying-samples.git")
	if err != nil {
		t.Error(err)
	}
	t.Log(buf.String())
}
func TestDockerlizeNpm(t *testing.T) {
	root, err := utils.Clone("https://github.com/jysperm/deploying-samples.git", "npm")
	buf, err := Dockerlize(root, nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(buf.String())
}

func TestDockerlizeYarn(t *testing.T) {
	root, err := utils.Clone("https://github.com/jysperm/deploying-samples.git", "yarn")
	buf, err := Dockerlize(root, nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(buf.String())
}
