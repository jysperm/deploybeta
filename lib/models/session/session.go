package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	accountModel "github.com/jysperm/deploying/lib/models/account"
	"github.com/jysperm/deploying/lib/services/etcd"
)

type Session struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

func CreateToken(account *accountModel.Account) (*Session, error) {
	sessionToken, err := randomString(32)

	if err != nil {
		return nil, err
	}

	sessionKey := fmt.Sprint("/sessions/", sessionToken)

	session := &Session{
		Token:    sessionToken,
		Username: account.Username,
	}

	jsonBytes, err := json.Marshal(account)

	if err != nil {
		return nil, err
	}

	_, err = etcd.Keys.Create(context.Background(), sessionKey, string(jsonBytes))

	return session, err
}

func randomString(length int) (string, error) {
	buffer := make([]byte, length)

	_, err := rand.Read(buffer)

	if err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(buffer)

	return strings.Replace(base64String, "/", "-", -1), nil
}
