package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jysperm/deploying/lib/etcd"
	"github.com/jysperm/deploying/lib/utils"
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

func FindSessionByToken(token string) (*Session, error) {
	sessionKey := fmt.Sprint("/sessions/", token)

	resp, err := etcd.Client.Get(context.Background(), sessionKey)

	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, ErrTokenNotFound
	}

	session := &Session{}

	err = json.Unmarshal([]byte(resp.Kvs[0].Value), session)

	if err != nil {
		return nil, err
	}

	return session, nil
}

func DeleteSessionByToken(token string) error {
	sessionKey := fmt.Sprint("/sessions/", token)

	_, err := etcd.Client.Delete(context.Background(), sessionKey)

	return err
}
