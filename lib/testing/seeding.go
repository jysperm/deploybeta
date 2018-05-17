package testing

import (
	"strings"

	"github.com/jysperm/deploybeta/lib/models"
	"github.com/jysperm/deploybeta/lib/utils"
)

func SeedAccount() (account models.Account, password string) {
	account = models.Account{
		Username: utils.RandomString(10),
		Email:    utils.RandomString(10) + "@gmail.com",
	}

	password = utils.RandomString(10)

	err := models.RegisterAccount(&account, password)

	if err != nil {
		panic(err)
	}

	return account, password
}

func SeedSession(account *models.Account) models.Session {
	session, err := models.CreateSession(account)

	if err != nil {
		panic(err)
	}

	return *session
}

func SeedApp(gitRepository string, owner string) models.Application {
	app := models.Application{
		Name:          strings.ToLower(utils.RandomString(10)),
		GitRepository: gitRepository,
		Instances:     1,
		Owner:         owner,
	}

	if err := models.CreateApp(&app); err != nil {
		panic(err)
	}

	return app
}
