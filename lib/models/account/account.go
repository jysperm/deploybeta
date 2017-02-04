package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	etcdClient "github.com/coreos/etcd/client"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"

	"github.com/jysperm/deploying/lib/services/etcd"
)

var ErrInvalidUsername = errors.New("invalid username")

type Account struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	Email        string `json:"email"`
}

var validUsername = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func Register(account *Account, password string) error {
	if !validUsername.MatchString(account.Username) {
		return ErrInvalidUsername
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	account.PasswordHash = string(passwordHash)

	accountKey := fmt.Sprint("/accounts/", account.Username)
	jsonBytes, err := json.Marshal(account)

	if err != nil {
		return err
	}

	_, err = etcd.Keys.Create(context.Background(), accountKey, string(jsonBytes))

	return err
}

func FindByName(username string) (*Account, error) {
	accountKey := fmt.Sprint("/accounts/", username)

	resp, err := etcd.Keys.Get(context.Background(), accountKey, &etcdClient.GetOptions{})

	if err != nil {
		return nil, err
	}

	account := &Account{}

	err = json.Unmarshal([]byte(resp.Node.Value), account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (account *Account) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
}
