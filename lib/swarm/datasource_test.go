package swarm

import (
	"testing"

	"github.com/jysperm/deploybeta/lib/models"
)

func TestCreateAndUpdateDataSource(t *testing.T) {
	dataSource := &models.DataSource{
		Name:          "test-redis",
		Type:          "redis",
		Instances:     2,
		OwnerUsername: "",
	}

	if err := UpdateDataSource(dataSource); err != nil {
		t.Fatal(err)
	}

	dataSource.Instances = 3

	if err := UpdateDataSource(dataSource); err != nil {
		t.Fatal(err)
	}

	if err := RemoveDataSource(dataSource); err != nil {
		t.Fatal(err)
	}
}
