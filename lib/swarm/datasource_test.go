package swarm

import (
	"context"
	"testing"
	"time"

	"github.com/jysperm/deploying/lib/models"
)

func TestCreateDataSource(t *testing.T) {
	datasource := models.DataSource{
		Name:      "test-redis-1",
		Type:      "redis",
		Instances: 2,
		Owner:     "",
	}

	if err := UpdateDataSource(&datasource, 2); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	if err := cleanup(datasource); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateDataSource(t *testing.T) {
	datasource := models.DataSource{
		Name:      "test-mongo",
		Type:      "mongodb",
		Instances: 2,
		Owner:     "",
	}

	if err := UpdateDataSource(&datasource, 2); err != nil {
		t.Fatal(err)
	}

	if err := UpdateDataSource(&datasource, 1); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)

	if err := cleanup(datasource); err != nil {
		t.Fatal(err)
	}

}

func TestRemoveDataSource(t *testing.T) {
	datasource := models.DataSource{
		Name:      "test-redis-none",
		Type:      "redis",
		Instances: 2,
		Owner:     "",
	}

	if err := UpdateDataSource(&datasource, 2); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	if err := RemoveDataSource(&datasource); err != nil {
		t.Fatal(err)
	}
}
func cleanup(datasource models.DataSource) error {
	serviceID, err := RetrieveServiceID(datasource.SwarmServiceName())
	err = swarmClient.ServiceRemove(context.Background(), serviceID)
	if err != nil {
		return err
	}

	err = RemoveOverlay(&datasource)
	if err != nil {
		return err
	}
	return nil
}
