package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"

	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/lib/db"
)

var ErrVersionNotFound = errors.New("version not found")

type Version struct {
	db.ResourceMeta

	AppName  string `json:"appName"`
	Tag      string `json:"tag"`
	Registry string `json:"registry"`
	Status   string `json:"status"`
}

func (version *Version) ResourceKey() string {
	return fmt.Sprintf("/apps/%s/versions/%s", version.AppName, version.Tag)
}

func (version *Version) Associations() []db.Association {
	return []db.Association{
		version.App(),
	}
}

func (version *Version) App() db.BelongsToAssociation {
	return db.BelongsTo((&Application{Name: version.AppName}).ResourceKey())
}

func NewVersion(app *Application) Version {
	return Version{
		AppName:  app.Name,
		Tag:      newVersionTag(),
		Registry: config.DefaultRegistry,
		Status:   "building",
	}
}

func FindVersionByTag(app *Application, tag string) (*Version, error) {
	version := &Version{
		AppName: app.Name,
		Tag:     tag,
	}

	err := db.Fetch(version)

	if err == db.ErrResourceNotFound {
		return nil, errwrap.Wrap(ErrVersionNotFound, err)
	}

	return version, err
}

func (version *Version) Create() error {
	_, err := db.StartTransaction(func(tran *db.Transaction) {
		tran.Create(version)
	})

	if err != nil {
		return errwrap.Wrapf("create version: {{err}}", err)
	}

	return nil
}

func (version *Version) Destroy() error {
	_, err := db.StartTransaction(func(tran *db.Transaction) {
		tran.Delete(version)
	})

	return err
}

func (version *Version) UpdateStatus(app *Application, status string) error {
	_, err := db.StartTransaction(func(tran *db.Transaction) {
		err := db.Fetch(version)

		if err != nil {
			tran.SetError(err)
			return
		}

		version.Status = status

		tran.Update(version)
	})

	return err
}

func (version *Version) DockerImageName() string {
	return fmt.Sprintf("%s/%s%s:%s", version.Registry, config.DockerPrefix, version.AppName, version.Tag)
}

func newVersionTag() string {
	now := time.Now()
	return fmt.Sprintf("%02d%02d%02d-%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}
