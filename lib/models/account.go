package models

import (
  "fmt"
  "encoding/json"
  "golang.org/x/net/context"
)

import "github.com/jysperm/deploying/lib/services/etcd"

type Account struct {
  Username string `json:"username"`
  Password string `json:"-"`
  PasswordHash string `json:"passwordHash"`
  Email string `json:"email"`
}

func CreateAccount(account *Account) *Account {
  key := fmt.Sprint("/accounts/", account.Username)
  value, err := json.Marshal(account)

  if err != nil {
    panic(err)
  }

  resp, err := etcd.Keys.Create(context.Background(), key, string(value))

  if err != nil {
    panic(err)
  }

  fmt.Println(resp)

  return account
}
