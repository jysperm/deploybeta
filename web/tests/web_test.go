package tests

import (
	"log"
	"os"
	"testing"

	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/lib/models"
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
	globalApp = SeedApp("https://github.com/jysperm/deploybeta-samples.git", globalAccount.Username)

	exitCode := m.Run()

	if err := globalSession.Destroy(); err != nil {
		log.Println(err)
	}

	if err := globalAccount.Destroy(); err != nil {
		log.Println(err)
	}

	if err := globalApp.Destroy(); err != nil {
		log.Println(err)
	}

	os.Exit(exitCode)
}
