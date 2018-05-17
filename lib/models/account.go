package models

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/jysperm/deploybeta/lib/etcd"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
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

func RegisterAccount(account *Account, password string) error {
	if !validUsername.MatchString(account.Username) {
		return ErrInvalidUsername
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	account.PasswordHash = string(passwordHash)

	accountKey := fmt.Sprint("/accounts/", account.Username)

	tran := etcd.NewTransaction()

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

func FindAccountByName(username string) (account Account, err error) {
	found, err := etcd.LoadKey(fmt.Sprintf("/accounts/%s", username), &account)

	if err != nil {
		return account, err
	} else if !found {
		return account, ErrAccountNotFound
	} else {
		return account, nil
	}
}

func DeleteAccountByName(username string) error {
	accountKey := fmt.Sprint("/accounts/", username)

	_, err := etcd.Client.Delete(context.Background(), accountKey)

	return err
}

func (account *Account) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
}
