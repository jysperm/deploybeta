package models

import (
	"errors"
	"fmt"

	"github.com/hashicorp/errwrap"

	"github.com/jysperm/deploybeta/lib/db"
	"github.com/jysperm/deploybeta/lib/utils"
)

var ErrTokenConflict = errors.New("token conflict")
var ErrTokenNotFound = errors.New("token not found")

type Session struct {
	db.ResourceMeta

	Token    string `json:"token"`
	Username string `json:"username"`
}

func (session *Session) ResourceKey() string {
	return fmt.Sprintf("/sessions/%s", session.Token)
}

func (session *Session) Associations() []db.Association {
	return []db.Association{}
}

func (session *Session) Account() db.HasOneAssociation {
	return db.HasOne((&Account{Username: session.Username}).ResourceKey())
}

func CreateSession(account *Account) (*Session, error) {
	session := &Session{
		Token:    utils.RandomString(32),
		Username: account.Username,
	}

	_, err := db.StartTransaction(func(tran *db.Transaction) {
		tran.Create(session)
	})

	return session, err
}

func FindSessionByToken(token string) (*Session, error) {
	session := &Session{
		Token: token,
	}

	err := db.Fetch(session)

	if err == db.ErrResourceNotFound {
		return nil, errwrap.Wrap(ErrTokenNotFound, err)
	}

	return session, err
}

func (session *Session) Destroy() error {
	_, err := db.StartTransaction(func(tran *db.Transaction) {
		tran.Delete(session)
	})

	return err
}
