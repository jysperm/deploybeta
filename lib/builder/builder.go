package builder

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jysperm/deploying/lib/etcd"

	"github.com/jysperm/deploying/config"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"golang.org/x/net/context"

	"github.com/jysperm/deploying/lib/models"
	"github.com/jysperm/deploying/lib/runtimes"
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

func BuildVersion(app *models.Application, gitTag string) (*models.Version, error) {
	versionTag := newTag()

	nameVersion := fmt.Sprintf("%s/%s:%s", config.DefaultRegistry, app.Name, versionTag)
	buildOpts := types.ImageBuildOptions{
		Tags:           []string{nameVersion},
		Dockerfile:     "Dockerfile",
		NoCache:        false,
		Remove:         true,
		SuppressOutput: false,
	}

	dirPath, err := cloneRepository(app.GitRepository, gitTag)
	if err != nil {
		return nil, err
	}

	fileBuffer, err := runtimes.Dockerlize(dirPath, app.GitRepository)
	if err != nil {
		return nil, err
	}

	if err := writeDockerfile(dirPath, fileBuffer); err != nil {
		return nil, err
	}

	buildCtx, err := buildContext(dirPath)
	if err != nil {
		return nil, err
	}

	defer buildCtx.Close()
	defer os.RemoveAll(dirPath)

	res, err := swarmClient.ImageBuild(context.Background(), buildCtx, buildOpts)
	if err != nil {
		return nil, err
	}

	v, err := models.CreateVersion(app, gitTag, versionTag)
	if err != nil {
		return nil, err
	}

	go wrtieProgress(app, versionTag, res.Body)

	return v, nil
}

func wrtieProgress(app *models.Application, tag string, r io.ReadCloser) {
	defer r.Close()

	reader := bufio.NewReader(r)
	v, err := models.FindVersionByTag(app, tag)
	if err != nil {
		errorKey := fmt.Sprintf("/apps/%s/versions/%s/progress/%d", app.Name, tag, time.Now().UnixNano())
		etcd.Client.Put(context.Background(), errorKey, err.Error())
	}

	for {
		eventKey := fmt.Sprintf("/apps/%s/versions/%s/progress/%d", app.Name, tag, time.Now().UnixNano())
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		_, err = etcd.Client.Put(context.Background(), eventKey, string(line))
		if err != nil {
			errorKey := fmt.Sprintf("/apps/%s/versions/%s/progress/%d", app.Name, tag, time.Now().UnixNano())
			etcd.Client.Put(context.Background(), errorKey, err.Error())
		}
		if strings.Contains(string(line), "errorDetail") {
			v.UpdateStatus(app, "fail")
			eventKey := fmt.Sprintf("/apps/%s/versions/%s/progress/%d", app.Name, tag, time.Now().UnixNano())
			etcd.Client.Put(context.Background(), eventKey, "Deploying: Building finished.")
			return
		}
	}

	nameVersion := fmt.Sprintf("%s/%s:%s", config.DefaultRegistry, app.Name, tag)

	if err := pushVersion(nameVersion); err != nil {
		errorKey := fmt.Sprintf("/apps/%s/versions/%s/progress/%d", app.Name, tag, time.Now().UnixNano())
		etcd.Client.Put(context.Background(), errorKey, err.Error())
		v.UpdateStatus(app, "fail")
	}

	v.UpdateStatus(app, "success")

	eventKey := fmt.Sprintf("/apps/%s/versions/%s/progress/%d", app.Name, tag, time.Now().UnixNano())
	if _, err := etcd.Client.Put(context.Background(), eventKey, "Deploying: Building finished."); err != nil {
		errorKey := fmt.Sprintf("/apps/%s/versions/%s/progress/%d", app.Name, tag, time.Now().UnixNano())
		etcd.Client.Put(context.Background(), errorKey, err.Error())
	}
}
func pushVersion(name string) error {
	res, err := swarmClient.ImagePush(context.Background(), name, types.ImagePushOptions{All: true, RegistryAuth: RegistryAuthParam})
	if err != nil {
		return err
	}

	reader := bufio.NewReader(res)
	for {
		_, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func newTag() string {
	now := time.Now()
	return fmt.Sprintf("%d%d%d-%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func cloneRepository(url string, gitTag string) (string, error) {
	if gitTag == "" {
		gitTag = "master"
	}

	path, err := utils.Clone(url, gitTag)
	if err != nil {
		return "", err
	}

	return path, nil
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

func buildContext(path string) (io.ReadCloser, error) {
	content, err := archive.Tar(path, archive.Gzip)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func writeDockerfile(path string, buf *bytes.Buffer) error {
	dockerfilePath := filepath.Join(path, "Dockerfile")
	Dockerfile, err := os.OpenFile(dockerfilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_TRUNC, 0666)
	defer Dockerfile.Close()
	if err != nil {
		return err
	}

	_, err = Dockerfile.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil

}
