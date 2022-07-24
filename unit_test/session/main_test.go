package session

import (
	"wellnus/backend/config"
	"wellnus/backend/db"
	"wellnus/backend/db/model"
	"wellnus/backend/router/session"
	"wellnus/backend/unit_test/test_helper"

	"testing"
	"os"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"database/sql"
	_ "github.com/lib/pq"
)

type User = model.User
type SessionResponse = model.SessionResponse

var (
	DB *sql.DB
	Router *gin.Engine
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
	
	DB = db.ConnectTestDB()
	Router = setupRouter()
	test_helper.ResetDB(DB)
	var err error

	testUsers, err = test_helper.SetupUsers(DB, 1)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err)) }
	
	os.Exit(m.Run())
}