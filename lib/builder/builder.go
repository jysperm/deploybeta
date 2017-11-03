package builder

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jysperm/deploying/config"
	"github.com/jysperm/deploying/lib/etcd"

	etcdv3 "github.com/coreos/etcd/clientv3"
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
var defaultTTL int64 = 60 * 10

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

	v, err := models.CreateVersion(app, versionTag)
	if err != nil {
		return nil, err
	}

	go wrtieProgress(app, versionTag, res.Body)

	return v, nil
}

func wrtieEvent(app *models.Application, lease *etcdv3.LeaseGrantResponse, tag string, event string) error {
	eventKey := fmt.Sprintf("/apps/%s/version/%s/progress/%s", app.Name, tag, time.Now().UnixNano())
	if _, err := etcd.Client.Put(context.Background(), eventKey, event, etcdv3.WithLease(lease.ID)); err != nil {
		return err
	}
	return nil
}

func wrtieProgress(app *models.Application, tag string, r io.ReadCloser) {
	defer r.Close()

	ttl, err := etcd.Client.Lease.Grant(context.Background(), defaultTTL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	reader := bufio.NewReader(r)
	v, err := models.FindVersionByTag(app, tag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	for {
		if _, err := etcd.Client.KeepAlive(context.Background(), ttl.ID); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			fmt.Fprintln(os.Stderr, err)
		}

		if err := wrtieEvent(app, ttl, tag, string(line)); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if strings.Contains(string(line), "errorDetail") {
			v.UpdateStatus(app, "fail")
			if err := wrtieEvent(app, ttl, tag, "Deploying: Building finished."); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		}
	}

	nameVersion := fmt.Sprintf("%s/%s:%s", config.DefaultRegistry, app.Name, tag)

	if err := pushVersion(nameVersion); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if err := v.UpdateStatus(app, "fail"); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	if err := v.UpdateStatus(app, "success"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if err := wrtieEvent(app, ttl, tag, "Deploying: Building finished."); err != nil {
		fmt.Fprintln(os.Stderr, err)
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
