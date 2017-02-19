package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	etcd "github.com/coreos/etcd/clientv3"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	"github.com/jysperm/deploying/lib/services"
	"github.com/jysperm/deploying/lib/utils"
)

var ErrTokenConflict = errors.New("token conflict")
var ErrTokenNotFound = errors.New("token not found")

type Session struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

func CreateToken(account *accountModel.Account) (*Session, error) {
	sessionToken := utils.RandomString(32)
	sessionKey := fmt.Sprint("/sessions/", sessionToken)

	session := &Session{
		Token:    sessionToken,
		Username: account.Username,
	}

	jsonBytes, err := json.Marshal(session)

	if err != nil {
		return nil, err
	}

	resp, err := services.EtcdClient.Txn(context.Background()).
		If(etcd.CreateRevision(sessionKey)).
		Then(etcd.OpPut(sessionKey, string(jsonBytes))).
		Commit()

	if err != nil {
		return nil, err
	}

	if resp.Succeeded == false {
		return nil, ErrTokenConflict
	}

	return session, err
}

func FindByToken(token string) (*Session, error) {
	sessionKey := fmt.Sprint("/sessions/", token)

	resp, err := services.EtcdClient.Get(context.Background(), sessionKey)

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

func DeleteByToken(token string) error {
	sessionKey := fmt.Sprint("/sessions/", token)

	_, err := services.EtcdClient.Delete(context.Background(), sessionKey)

	return err
}
