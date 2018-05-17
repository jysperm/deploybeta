package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/lib/models"
	"github.com/jysperm/deploybeta/lib/swarm"
	. "github.com/jysperm/deploybeta/lib/testing"
	"github.com/jysperm/deploybeta/web"
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
