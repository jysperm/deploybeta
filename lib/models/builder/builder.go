package builder

import (
	"crypto/sha1"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"golang.org/x/net/context"
	"srcd.works/go-git.v4"
)

func pullRepository(url string, depth int) (string, error) {
	tempDir := os.TempDir()
	hash := sha1.New()
	io.WriteString(hash, url)
	newDir := tempDir + string(hash.Sum(nil))
	err := os.MkdirAll(newDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	_, err = git.PlainClone(newDir, false, &git.CloneOptions{
		URL:   url,
		Depth: depth,
	})
	if err != nil {
		return "", err
	}

	return newDir, nil
}

func tarRepository(path string) (io.ReadCloser, error) {
	content, err := archive.Tar(path, archive.Gzip)
	if err != nil {
		return nil, err
	}

	return content, nil

}

//BuildImage will build a docker image accroding to the repo's url and depth and Dockerfiles
func BuildImage(opts types.ImageBuildOptions, url string, depth int) error {
	client, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	dirPath, err := pullRepository(url, depth)
	if err != nil {
		return err
	}

	content, err := tarRepository(dirPath)
	if err != nil {
		return err
	}

	response, err := client.ImageBuild(ctx, content, opts)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}
