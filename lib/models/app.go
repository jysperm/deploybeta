package models

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/errwrap"

	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/lib/db"
)

var ErrInvalidName = errors.New("invalid app name")
var ErrUpdateConflict = errors.New("update version conflict")
var ErrAppNotFound = errors.New("app not found")

type Application struct {
	db.ResourceMeta

	Name          string `json:"name"`
	OwnerUsername string `json:"ownerUsername"`
	GitRepository string `json:"gitRepository"`
	Instances     int    `json:"instances"`
	VersionTag    string `json:"versionTag"`
}

func (app *Application) ResourceKey() string {
	return fmt.Sprintf("/apps/%s", app.Name)
}

func (app *Application) Associations() []db.Association {
	return []db.Association{
		app.Owner(),
		app.Upstreams(),
		app.Version(),
		app.Versions(),
		app.DataSources(),
	}
}

func (app *Application) Owner() db.HasOneAssociation {
	return db.BelongsTo(
		(&Account{Username: app.OwnerUsername}).ResourceKey(),
		fmt.Sprintf("/accounts/%s/apps", app.OwnerUsername),
	)
}

func (app *Application) Upstreams() db.HasManyAssociation {
	return db.HasManyThrough(fmt.Sprintf("/apps/%s/upstreams", app.Name))
}

func (app *Application) Version() db.HasOneAssociation {
	return db.HasOne(fmt.Sprintf("/apps/%s/versions/%s", app.Name, app.VersionTag))
}

func (app *Application) Versions() db.HasManyAssociation {
	return db.HasManyPrefix(fmt.Sprintf("/apps/%s/versions/", app.Name))
}

func (app *Application) DataSources() db.HasManyAssociation {
	return db.HasManyThrough(fmt.Sprintf("/apps/%s/data-sources", app.Name))
}

var validName = regexp.MustCompile(`^[a-z0-9_-]+$`)

func CreateApp(app *Application) error {
	if !validName.MatchString(app.Name) {
		return ErrInvalidName
	}

	_, err := db.StartTransaction(func(tran db.Transaction) {
		tran.Create(app)

		app.Upstreams().Attach(tran, tran.Create(&Upstream{
			Domain:   fmt.Sprint(app.Name, config.WildcardDomain),
			Backends: []UpstreamBackend{},
		}))
	})

	return err
}

func FindAppByName(name string) (*Application, error) {
	app := &Application{
		Name: name,
	}

	err := db.Fetch(app)

	if err == db.ErrResourceNotFound {
		return nil, errwrap.Wrap(ErrAppNotFound, err)
	}

	return app, err
}

func (app *Application) Destroy() error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		tran.Delete(app)
	})

	return err
}

func (app *Application) Update(updates *Application) error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		err := db.Fetch(app)

		if err != nil {
			tran.SetError(err)
			return
		}

		if updates.GitRepository != "" {
			app.GitRepository = updates.GitRepository
		}

		if updates.Instances != 0 {
			app.Instances = updates.Instances
		}

		tran.Update(app)
	})

	return err
}

func (app *Application) AddUpstream(domain string) error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		app.Upstreams().Attach(tran, tran.Create(&Upstream{
			Domain:   domain,
			Backends: []UpstreamBackend{},
		}))
	})

	return err
}

func (app *Application) RemoveUpstream(domain string) error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		app.Upstreams().Detach(tran, tran.Remove(&Upstream{
			Domain: domain,
		}))
	})

	return err
}

func (app *Application) UpdateVersion(versionTag string) error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		err := db.Fetch(app)

		if err != nil {
			tran.SetError(err)
			return
		}

		app.VersionTag = versionTag

		tran.Update(app)
	})

	return err
}

func (app *Application) SwarmServiceName() string {
	return fmt.Sprintf("%s%s", config.DockerPrefix, app.Name)
}

func (app *Application) SwarmInstances() int {
	return app.Instances
}
