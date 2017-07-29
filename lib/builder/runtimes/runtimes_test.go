package runtimes

import (
	"testing"

	"github.com/jysperm/deploying/lib/utils"
)

func TestDockerlizeNode(t *testing.T) {
	root, err := utils.Clone("https://github.com/jysperm/deploying-samples.git", "npm")
	err = Dockerlize(root, nil)
	if err != nil {
		t.Error(err)
	}
}
