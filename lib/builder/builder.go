package builder

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	etcdv3 "github.com/coreos/etcd/clientv3"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/hashicorp/errwrap"
	"golang.org/x/net/context"

	"github.com/jysperm/deploybeta/lib/db"
	"github.com/jysperm/deploybeta/lib/models"
	"github.com/jysperm/deploybeta/lib/runtimes"
	"github.com/jysperm/deploybeta/lib/utils"
)

const RegistryAuthParam = "deploybeta"

var swarmClient *client.Client
var defaultTTL int64 = 60 * 10

func init() {
	var err error
	swarmClient, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

type buildEvent struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

func BuildVersion(app *models.Application, gitTag string) (*models.Version, error) {
	fmt.Printf("Clone from %s\n", app.GitRepository)

	dirPath, err := cloneRepository(app.GitRepository, gitTag)
	if err != nil {
		return nil, errwrap.Wrapf("clone repository: {{err}}", err)
	}

	runtime, err := runtimes.DecideRuntime(runtimes.NewBuildContext(dirPath, app.GitRepository))
	if err != nil {
		return nil, errwrap.Wrapf("decide runtime: {{err}}", err)
	}

	fileBuffer, err := runtime.Dockerfile()
	if err != nil {
		return nil, errwrap.Wrapf("generate dockerfile: {{err}}", err)
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

	version := models.NewVersion(app)

	fmt.Printf("Building image %s\n", version.DockerImageName())

	res, err := swarmClient.ImageBuild(context.Background(), buildCtx, types.ImageBuildOptions{
		Tags:           []string{version.DockerImageName()},
		Dockerfile:     "Dockerfile",
		NoCache:        false,
		Remove:         true,
		SuppressOutput: false,
	})

	if err != nil {
		return nil, err
	}

	err = version.Create()
	if err != nil {
		return nil, err
	}

	go wrtieProgress(app, &version, res.Body)

	return &version, nil
}

func wrtieEvent(app *models.Application, lease *etcdv3.LeaseGrantResponse, tag string, event string) error {
	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	newEvent := buildEvent{
		ID:      id,
		Payload: event,
	}
	eventKey := fmt.Sprintf("/progress/%s/%s/%s", app.Name, tag, id)
	e, err := json.Marshal(newEvent)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if _, err := db.Client.Put(context.Background(), eventKey, string(e), etcdv3.WithLease(lease.ID)); err != nil {
		return err
	}
	return nil
}

func wrtieProgress(app *models.Application, version *models.Version, r io.ReadCloser) {
	defer r.Close()

	ttl, err := db.Client.Lease.Grant(context.Background(), defaultTTL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	reader := bufio.NewReader(r)

	for {
		if _, err := db.Client.KeepAlive(context.Background(), ttl.ID); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			fmt.Fprintln(os.Stderr, err)
		}

		if err := wrtieEvent(app, ttl, version.Tag, string(line)); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if strings.Contains(string(line), "errorDetail") {
			version.UpdateStatus(app, "fail")
			if err := wrtieEvent(app, ttl, version.Tag, "Deploybeta: Building finished."); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		}
	}

	if err := pushVersion(version); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if err := version.UpdateStatus(app, "fail"); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	if err := version.UpdateStatus(app, "success"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if err := wrtieEvent(app, ttl, version.Tag, "Deploybeta: Building finished."); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func pushVersion(version *models.Version) error {
	res, err := swarmClient.ImagePush(context.Background(), version.DockerImageName(), types.ImagePushOptions{All: true, RegistryAuth: RegistryAuthParam})
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
	Dockerfile, err := os.OpenFile(dockerfilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
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
