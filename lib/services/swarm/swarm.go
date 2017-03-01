package swarm

import (
	dockerSwarm "github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

var swarmClient *client.Client

func init() {
	var err error
	swarmClient, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

//InitSwarm will initialize a swarm and return the id of initialized swarm
// and return its id
func InitSwarm(req dockerSwarm.InitRequest) (string, error) {
	ctx := context.Background()

	id, err := swarmClient.SwarmInit(ctx, req)
	if err != nil {
		return "", err
	}
	return id, nil
}

//LeaveSwarm will leave the current swarm
func LeaveSwarm(force bool) error {
	ctx := context.Background()

	if err := swarmClient.SwarmLeave(ctx, force); err != nil {
		return err
	}
	return nil
}

//JoinSwarm will join a swarm
func JoinSwarm(req dockerSwarm.JoinRequest) error {
	ctx := context.Background()

	if err := swarmClient.SwarmJoin(ctx, req); err != nil {
		return err
	}

	return nil
}

//UpdateSwarm will update the config of swarm
func UpdateSwarm(version dockerSwarm.Version, spec dockerSwarm.Spec, flags dockerSwarm.UpdateFlags) error {
	ctx := context.Background()

	if err := swarmClient.SwarmUpdate(ctx, version, spec, flags); err != nil {
		return err
	}

	return nil
}
