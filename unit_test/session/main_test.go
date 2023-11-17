package session

import (
	"wellnus/backend/config"
	"wellnus/backend/db"
	. "wellnus/backend/db/model"
	"wellnus/backend/router/session"
	"wellnus/backend/unit_test/test_helper"

	"fmt"
	"log"
	"os"
	"testing"

	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	DB          *sql.DB
	Router      *gin.Engine
	sessionKey1 string
	sessionKey2 string
)

var testUsers []User
var sessionKey string

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/session", session.LoginHandler(DB))
	router.DELETE("/session", session.LogoutHandler(DB))

	return router
}

func TestMain(m *testing.M) {
	config.LoadENV("../../.env")

	DB = db.ConnectDB()
	Router = setupRouter()
	test_helper.ResetDB(DB)
	var err error

	testUsers, err = test_helper.SetupUsers(DB, 1)
	if err != nil {
		log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err))
	}

	os.Exit(m.Run())
}
