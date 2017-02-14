package builder

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	dockerbuilder "github.com/docker/docker/builder"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/gitutils"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func pullRepository(url string) (string, error) {
	path, err := gitutils.Clone(url)
	if err != nil {
		return "", err
	}
	return path, nil
}

func tarRepository(path string) (io.ReadCloser, error) {
	content, err := archive.Tar(path, archive.Gzip)
	if err != nil {
		return nil, err
	}

	return content, nil
}

//BuildImage will build a docker image accroding to the repo's url and depth and Dockerfiles
func BuildImage(opts types.ImageBuildOptions, url string) error {
	switch {
	case opts.Tags == nil:
		return errors.New("Need one or more tags")
	case opts.Dockerfile == "":
		return errors.New("Need a name of Dockerfile")
	}
	opts.NoCache = true
	opts.Remove = true
	opts.ForceRemove = true
	opts.SuppressOutput = true
	opts.Isolation = ""

	client, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	dirPath, err := pullRepository(url)
	if err != nil {
		return err
	}

	content, err := tarRepository(dirPath)
	defer content.Close()
	if err != nil {
		return err
	}

	buildCtx, relDockerfile, err := dockerbuilder.GetContextFromReader(content, opts.Dockerfile)
	if err != nil {
		return err
	}
	opts.Dockerfile = relDockerfile
	response, err := client.ImageBuild(ctx, buildCtx, opts)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	defer os.RemoveAll(dirPath)

	return nil
}
