package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/jysperm/deploybeta/lib/etcd"
	"github.com/jysperm/deploybeta/lib/utils"
)

var ErrTokenConflict = errors.New("token conflict")
var ErrTokenNotFound = errors.New("token not found")

type Session struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

func CreateSession(account *Account) (*Session, error) {
	sessionToken := utils.RandomString(32)
	sessionKey := fmt.Sprint("/sessions/", sessionToken)

	session := &Session{
		Token:    sessionToken,
		Username: account.Username,
	}

	tran := etcd.NewTransaction()

	tran.CreateJSON(sessionKey, session)

	resp, err := tran.Execute()

	if err != nil {
		return nil, err
	}

	if resp.Succeeded == false {
		return nil, ErrTokenConflict
	}

	return session, err
}

func FindSessionByToken(token string) (session Session, err error) {
	found, err := etcd.LoadKey(fmt.Sprintf("/sessions/%s", token), &session)

	if err != nil {
		return session, err
	} else if !found {
		return session, ErrTokenNotFound
	} else {
		return session, nil
	}
}

func DeleteSessionByToken(token string) error {
	sessionKey := fmt.Sprint("/sessions/", token)

	_, err := etcd.Client.Delete(context.Background(), sessionKey)

	return err
}
