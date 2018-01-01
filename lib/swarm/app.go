package swarm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jysperm/deploying/config"
	"github.com/jysperm/deploying/lib/etcd"
	"github.com/jysperm/deploying/lib/models"

	"github.com/docker/docker/api/types/swarm"
)

type UpstreamConfig struct {
	Port uint32 `json:"port"`
}

func UpdateApp(app *models.Application) error {
	currentVersion, err := models.FindVersionByTag(app, app.Version)
	if err != nil {
		return err
	}
	nameVersion := fmt.Sprintf("%s/%s:%s", config.DefaultRegistry, app.Name, currentVersion.Tag)

	var upstreamConfig UpstreamConfig

	if err := UpdateService(app.Name, uint64(app.Instances), []swarm.PortConfig{}, []swarm.NetworkAttachmentConfig{}, nameVersion); err != nil {
		return err
	}

	serviceID, _ := RetrieveServiceID(app.Name)
	upstreamConfig.Port, err = RetrievePort(serviceID)
	if err != nil {
		return err
	}

	upstream, err := json.Marshal([]UpstreamConfig{upstreamConfig})
	if err != nil {
		return err
	}

	upstreamKey := fmt.Sprintf("/upstream/%s", app.Name)
	if _, err := etcd.Client.Put(context.Background(), upstreamKey, string(upstream)); err != nil {
		return err
	}

	if err := app.Update(app); err != nil {
		return err
	}

	if err := app.UpdateVersion(currentVersion.Tag); err != nil {
		return err
	}

	return nil
}
