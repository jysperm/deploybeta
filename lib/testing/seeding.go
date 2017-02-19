package testing

import (
	accountModel "github.com/jysperm/deploying/lib/models/account"
	sessionModel "github.com/jysperm/deploying/lib/models/session"
	"github.com/jysperm/deploying/lib/utils"
)

func SeedAccount() (account accountModel.Account, password string) {
	account = accountModel.Account{
		Username: utils.RandomString(10),
		Email:    utils.RandomString(10) + "@gmail.com",
	}

	password = utils.RandomString(10)

	err := accountModel.Register(&account, password)

	if err != nil {
		panic(err)
	}

	return account, password
}

func SeedSession(account *accountModel.Account) sessionModel.Session {
	session, err := sessionModel.CreateToken(account)

	if err != nil {
		panic(err)
	}

	return *session
}
