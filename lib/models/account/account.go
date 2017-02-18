package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	etcd "github.com/coreos/etcd/clientv3"

	"github.com/jysperm/deploying/lib/services"
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

	resp, err := services.EtcdClient.Txn(context.Background()).
		If(etcd.Compare(etcd.CreateRevision(accountKey), "=", 0)).
		Then(etcd.OpPut(accountKey, string(jsonBytes))).
		Commit()

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
