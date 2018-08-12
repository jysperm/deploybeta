package models

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/errwrap"
	"github.com/jysperm/deploybeta/lib/db"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidUsername = errors.New("invalid username")
var ErrUsernameConflict = errors.New("username conflict")
var ErrAccountNotFound = errors.New("account not found")

var regexValidUsername = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

type Account struct {
	db.ResourceMeta

	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	Email        string `json:"email"`
}

func (account *Account) ResourceKey() string {
	return fmt.Sprintf("/accounts/%s", account.Username)
}

func (account *Account) Associations() []db.Association {
	return []db.Association{
		account.Apps(),
		account.DataSources(),
	}
}

func (account *Account) Apps() db.HasManyAssociation {
	return db.HasManyThrough(fmt.Sprintf("/accounts/%s/apps", account.Username))
}

func (account *Account) DataSources() db.HasManyAssociation {
	return db.HasManyThrough(fmt.Sprintf("/accounts/%s/data-sources", account.Username))
}

func RegisterAccount(account *Account, password string) error {
	if !regexValidUsername.MatchString(account.Username) {
		return ErrInvalidUsername
	}

	account.SetPassword(password)

	_, err := db.StartTransaction(func(tran db.Transaction) {
		tran.Create(account)
	})

	if err == db.ErrEtcdTransactionFailed {
		return errwrap.Wrap(ErrUsernameConflict, err)
	} else if err != nil {
		return err
	}

	return nil
}

func FindAccountByName(username string) (*Account, error) {
	account := &Account{
		Username: username,
	}

	err := db.Fetch(account)

	if err == db.ErrResourceNotFound {
		return nil, errwrap.Wrap(ErrAccountNotFound, err)
	}

	return account, err
}

func (account *Account) Destroy() error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		tran.Delete(account)
	})

	return err
}

func (account *Account) SetPassword(password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	account.PasswordHash = string(passwordHash)

	return nil
}

func (account *Account) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
}
