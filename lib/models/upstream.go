package models

import (
	"fmt"

	"github.com/jysperm/deploybeta/lib/db"
)

type Upstream struct {
	db.ResourceMeta

	Domain   string            `json:"domain"`
	Backends []UpstreamBackend `json:"backends"`
}

type UpstreamBackend struct {
	Port uint32 `json:"port"`
}

func (upstream *Upstream) ResourceKey() string {
	return fmt.Sprintf("/upstreams/%s", upstream.Domain)
}

func (upstream *Upstream) Associations() []db.Association {
	return []db.Association{}
}

func UpdateUpstreamForApp(app *Application, backends []UpstreamBackend) error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		upstreams := make([]Upstream, 0)
		err := app.Upstreams().FetchAll(&upstreams)

		if err != nil {
			tran.SetError(err)
			return
		}

		for _, upstream := range upstreams {
			upstream.Backends = backends
			tran.Put(&upstream)
		}
	})

	return err
}
