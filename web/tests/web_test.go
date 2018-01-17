package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/jysperm/deploying/config"
	"github.com/jysperm/deploying/lib/models"
	"github.com/jysperm/deploying/lib/swarm"
	. "github.com/jysperm/deploying/lib/testing"
	"github.com/jysperm/deploying/web"
)

func init() {
	go web.CreateWebServer().Start(config.Listen)
}

var globalAccount models.Account
var globalSession models.Session
var globalApp models.Application

func TestMain(m *testing.M) {
	globalAccount, _ = SeedAccount()
	globalSession = SeedSession(&globalAccount)
	globalApp = SeedApp("https://github.com/jysperm/deploying-samples.git", globalAccount.Username)

	fmt.Println(globalAccount)

	exitVal := m.Run()

	models.DeleteSessionByToken(globalSession.Token)
	models.DeleteAccountByName(globalAccount.Username)
	swarm.RemoveService(&globalApp)
	os.Exit(exitVal)
}
