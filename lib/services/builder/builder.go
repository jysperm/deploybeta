package builder

import (
	"crypto/sha1"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"golang.org/x/net/context"
	"srcd.works/go-git.v4"
	"srcd.works/go-git.v4/plumbing"
)

func pullRepository(url string, branch string) (string, error) {
	hash := sha1.New()
	io.WriteString(hash, url)
	hashURL := hash.Sum(nil)

	tempDir, err := ioutil.TempDir("", string(hashURL))
	if err != nil {
		return "", nil
	}

	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName(branch),
	})
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

func tarRepository(path string) (io.ReadCloser, error) {
	content, err := archive.Tar(path, archive.Gzip)
	if err != nil {
		return nil, err
	}

	return content, nil

}

//BuildImage will build a docker image accroding to the repo's url and depth and Dockerfiles
func BuildImage(opts types.ImageBuildOptions, url string, branch string) error {
	client, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	dirPath, err := pullRepository(url, branch)
	if err != nil {
		return err
	}

	content, err := tarRepository(dirPath)
	defer content.Close()
	if err != nil {
		return err
	}

	response, err := client.ImageBuild(ctx, content, opts)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	defer os.RemoveAll(dirPath)

	return nil
}
