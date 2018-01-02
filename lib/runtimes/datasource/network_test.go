package datasource

import (
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestCreateOverlay(t *testing.T) {
	id, err := CreateOverlay("relay")
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
	time.Sleep(1000 * time.Millisecond)
	if err := swarmClient.NetworkRemove(context.Background(), id); err != nil {
		t.Log(err)
	}

}

func TestRemoveOverlay(t *testing.T) {
	if _, err := CreateOverlay("relay"); err != nil {
		t.Error(err)
	}
	time.Sleep(1000 * time.Millisecond)
	if err := RemoveOverlay("relay"); err != nil {
		t.Error(err)
	}
	if id, err := FindByName("relay"); id != "" || err != nil {
		t.Error(err)
	}
}

func TestListOverlay(t *testing.T) {
	list, err := ListOverlays()
	if err != nil {
		t.Error(err)
	}
	t.Log(list)
}
