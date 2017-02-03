package account

import (
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"

	"github.com/jysperm/deploying/lib/services/etcd"
)

type Account struct {
	Username     string `json:"username"`
	Password     string `json:"-"`
	PasswordHash string `json:"passwordHash"`
	Email        string `json:"email"`
}

func Register(account *Account) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	account.PasswordHash = string(passwordHash)

	etcdKey := fmt.Sprint("/accounts/", account.Username)
	value, err := json.Marshal(account)

	if err != nil {
		return err
	}

	_, err = etcd.Keys.Create(context.Background(), etcdKey, string(value))

	return err
}
