package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"

	"github.com/jysperm/deploying/lib/services"
)

var ErrInvalidUsername = errors.New("invalid username")
var ErrUsernameConflict = errors.New("username conflict")
var ErrAccountNotFound = errors.New("account not found")

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

	tran := services.NewEtcdTransaction()

	tran.CreateJSON(accountKey, account)

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUsernameConflict
	}

	return nil
}

func FindByName(username string) (*Account, error) {
	accountKey := fmt.Sprint("/accounts/", username)

	resp, err := services.EtcdClient.Get(context.Background(), accountKey)

	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, ErrAccountNotFound
	}

	account := &Account{}

	err = json.Unmarshal([]byte(resp.Kvs[0].Value), account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func DeleteByName(username string) error {
	accountKey := fmt.Sprint("/accounts/", username)

	_, err := services.EtcdClient.Delete(context.Background(), accountKey)

	return err
}

func (account *Account) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
}
