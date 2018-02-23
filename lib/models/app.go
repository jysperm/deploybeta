package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/jysperm/deploying/config"
	"github.com/jysperm/deploying/lib/etcd"
	"golang.org/x/net/context"
)

var ErrInvalidName = errors.New("invalid app name")
var ErrUpdateConflict = errors.New("update version conflict")

type Application struct {
	Name          string `json:"name"`
	Owner         string `json:"owner"`
	GitRepository string `json:"gitRepository"`
	Instances     int    `json:"instances"`
	Version       string `json:"version"`
}

var validName = regexp.MustCompile(`^[a-z0-9_-]+$`)

func CreateApp(app *Application) error {
	if !validName.MatchString(app.Name) {
		return ErrInvalidName
	}

	tran := etcd.NewTransaction()

	tran.AppendStringArray(fmt.Sprintf("/account/%s/apps", app.Owner), app.Name)
	tran.CreateJSON(fmt.Sprint("/apps/", app.Name), app)

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUpdateConflict
	}

	return nil
}

// TODO: Delete app name from `/account/:name/apps`
func DeleteAppByName(name string) error {
	appKey := fmt.Sprint("/apps/", name)

	_, err := etcd.Client.Delete(context.Background(), appKey)

	return err
}

func GetAppsOfAccount(account *Account) (result []Application, err error) {
	accountAppsKey := fmt.Sprintf("/account/%s/apps", account.Username)
	resp, err := etcd.Client.Get(context.Background(), accountAppsKey)

	result = make([]Application, 0)

	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return result, nil
	}

	accountApps := []string{}
	err = json.Unmarshal([]byte(resp.Kvs[0].Value), &accountApps)

	if err != nil {
		return nil, err
	}

	for _, appName := range accountApps {
		appKey := fmt.Sprint("/apps/", appName)
		resp, err = etcd.Client.Get(context.Background(), appKey)

		if err != nil {
			return result, err
		}

		app := Application{}

		if len(resp.Kvs) != 0 {
			err = json.Unmarshal([]byte(resp.Kvs[0].Value), &app)

			if err != nil {
				return result, err
			}

			result = append(result, app)
		}
	}

	return result, nil
}

func (app *Application) Update(update *Application) error {
	appKey := fmt.Sprint("/apps/", app.Name)

	tran := etcd.NewTransaction()

	tran.WatchJSON(appKey, &Application{}, func(watchedKey interface{}) error {
		app := *watchedKey.(*Application)

		if update.GitRepository != "" {
			app.GitRepository = update.GitRepository
		}

		if update.Instances != 0 {
			app.Instances = update.Instances
		}

		update.Version = app.Version

		tran.PutJSON(appKey, app)

		return nil
	})

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUpdateConflict
	}

	if update.GitRepository != "" {
		app.GitRepository = update.GitRepository
	}

	if update.Instances != 0 {
		app.Instances = update.Instances
	}

	return nil
}

func (app *Application) UpdateVersion(version string) error {
	appKey := fmt.Sprint("/apps/", app.Name)

	tran := etcd.NewTransaction()

	tran.WatchJSON(appKey, &Application{}, func(watchedKey interface{}) error {
		app := *watchedKey.(*Application)

		app.Version = version

		tran.PutJSON(appKey, app)

		return nil
	})

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUpdateConflict
	}

	app.Version = version

	return nil
}

func (app *Application) SwarmServiceName() string {
	return fmt.Sprintf("%s%s", config.DockerPrefix, app.Name)
}

func FindAppByName(name string) (app Application, err error) {
	found, err := etcd.LoadKey(fmt.Sprintf("/apps/%s", name), &app)

	if err != nil {
		return app, err
	} else if !found {
		return app, errors.New("app not found")
	} else {
		return app, nil
	}
}
