package helpers

import accountModel "github.com/jysperm/deploying/lib/models/account"

type AccountResponse struct {
  Username string `json:"username"`
  Email string `json:"email"`
}

func NewAccountResponse(account *accountModel.Account) AccountResponse {
  return AccountResponse{
    Username: account.Username,
    Email: account.Email,
  }
}
