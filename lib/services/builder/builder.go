package builder

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/gitutils"
	"golang.org/x/net/context"
)

func cloneRepository(url string) (string, error) {
	path, err := gitutils.Clone(url)
	if err != nil {
		return "", err
	}
	return path, nil
}

func buildContext(path string) (io.ReadCloser, error) {
	content, err := archive.Tar(path, archive.Gzip)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func extractShasum(r io.ReadCloser) (string, error) {
	var shasum string
	var buildErr error
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		shasum, buildErr = func(s []byte) (string, error) {
			var f interface{}
			if err := json.Unmarshal(s, &f); err != nil {
				return "", err
			}
			m := f.(map[string]interface{})
			for k, v := range m {
				switch vv := v.(type) {
				case string:
					if k == "stream" && strings.HasPrefix(vv, "sha256") {
						return vv[len("sha256:") : len(vv)-1], nil
					}
					if k == "error" {
						return "", errors.New(vv)
					}
				}

			}
			return "", nil
		}(line)
		if buildErr != nil {
			return "", buildErr
		}
	}
	return shasum, nil
}

//BuildImage will build a docker image accroding to the repo's url and depth and Dockerfiles
func BuildImage(opts types.ImageBuildOptions, url string) (string, error) {
	if opts.Dockerfile == "" {
		opts.Dockerfile = "Dockerfile"
	}
	opts.NoCache = false
	opts.Remove = true
	opts.SuppressOutput = true
	opts.Isolation = ""

	client, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	dirPath, err := cloneRepository(url)
	if err != nil {
		return "", err
	}

	buildCtx, err := buildContext(dirPath)
	if err != nil {
		return "", err
	}
	defer buildCtx.Close()
	defer os.RemoveAll(dirPath)

	response, err := client.ImageBuild(ctx, buildCtx, opts)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	id, err := extractShasum(response.Body)
	if err != nil {
		return "", err
	}

	return id, nil
}
