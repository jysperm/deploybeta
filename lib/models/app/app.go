package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/jysperm/deploying/lib/etcd"
	accountModel "github.com/jysperm/deploying/lib/models/account"
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

	appKey := fmt.Sprint("/apps/", app.Name)
	accountAppsKey := fmt.Sprintf("/account/%s/apps", app.Owner)

	tran := etcd.NewTransaction()

	tran.WatchJSON(accountAppsKey, &[]string{})
	tran.CreateJSON(appKey, app)

	resp, err := tran.Execute(func(watchedKeys map[string]interface{}) error {
		accountApps := *watchedKeys[accountAppsKey].(*[]string)

		tran.PutJSONOnSuccess(accountAppsKey, append(accountApps, app.Name))

		return nil
	})

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUpdateConflict
	} else {
		return nil
	}
}

// TODO: Delete app name from `/account/:name/apps`
func DeleteByName(name string) error {
	appKey := fmt.Sprint("/apps/", name)

	_, err := etcd.Client.Delete(context.Background(), appKey)

	return err
}

func GetAppsOfAccount(account *accountModel.Account) (result []Application, err error) {
	accountAppsKey := fmt.Sprintf("/account/%s/apps", account.Username)
	resp, err := etcd.Client.Get(context.Background(), accountAppsKey)

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

func (app *Application) UpdateGitRepository(gitRepository string) error {
	return nil
}

func (app *Application) UpdateInstances(instances int) error {
	return nil
}

func (app *Application) UpdateVersion(version string) error {
	return nil
}

func FindByName(name string) (*Application, error) {
	appKey := fmt.Sprintf("/apps/%s", name)
	resp, err := etcd.Client.Get(context.Background(), appKey)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	var appFound Application
	if err := json.Unmarshal(resp.Kvs[0].Value, &appFound); err != nil {
		return nil, err
	}

	return &appFound, nil
}
