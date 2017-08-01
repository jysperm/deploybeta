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
	"golang.org/x/net/context"

	"github.com/jysperm/deploying/lib/builder/runtimes"
	"github.com/jysperm/deploying/lib/utils"
)

const RegistryAuthParam = "deploying"

var swarmClient *client.Client

func init() {
	var err error
	swarmClient, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

func cloneRepository(url string, param string) (string, error) {
	if param == "" {
		param = "master"
	}
	path, err := utils.Clone(url, param)
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

func PushImage(image string) error {
	if _, err := swarmClient.ImagePush(context.Background(), image, types.ImagePushOptions{All: true, RegistryAuth: RegistryAuthParam}); err != nil {
		return err
	}
	return nil
}

func LookupRepoTag(name string, id string) (string, error) {
	var tag string
	inspect, _, err := swarmClient.ImageInspectWithRaw(context.Background(), id)
	if err != nil {
		return "", err
	}

	for _, i := range inspect.RepoTags {
		if strings.Contains(i, name) {
			tag = i
			break
		}
	}
	return tag, nil
}

//BuildImage will build a docker image accroding to the repo's url and depth and Dockerfiles
func BuildImage(opts types.ImageBuildOptions, url string, param string) (string, error) {
	if opts.Dockerfile == "" {
		opts.Dockerfile = "Dockerfile"
	}
	opts.NoCache = false
	opts.Remove = true
	opts.SuppressOutput = true
	opts.Isolation = ""

	dirPath, err := cloneRepository(url, param)
	if err != nil {
		return "", err
	}

	if err := runtimes.Dockerlize(dirPath, url); err != nil {
		return "", err
	}

	buildCtx, err := buildContext(dirPath)
	if err != nil {
		return "", err
	}
	defer buildCtx.Close()
	defer os.RemoveAll(dirPath)

	response, err := swarmClient.ImageBuild(context.Background(), buildCtx, opts)
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
