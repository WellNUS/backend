package matching

import (
	"wellnus/backend/config"
	"wellnus/backend/db"
	"wellnus/backend/db/model"
	"wellnus/backend/router/match"
	"wellnus/backend/router/http_helper/http_error"
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
type MatchSetting = model.MatchSetting

var (
	DB *sql.DB 
	Router *gin.Engine
	NotFoundErrorMessage 		string = http_error.NotFoundError.Error()
	UnauthorizedErrorMessage	string = http_error.UnauthorizedError.Error()
)

var testUsers []User
var sessionKeys []string

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/match", match.AddMatchRequestHandler(DB))
	return router
}

func TestMain(m *testing.M) {
	config.LoadENV("../../.env")

	DB = db.ConnectDB()
	Router = setupRouter()
	
	test_helper.ResetDB(DB)
	var err error
	// Setup test users
	testUsers, err = test_helper.SetupUsers(DB, config.MATCH_THRESHOLD)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test users. %v", err)) }

	// Create match settings for all users
	_, err = test_helper.SetupMatchSettingForUsers(DB, testUsers)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating match settings. %v", err)) }

	// Create sessions for the first 2 users
	sessionKeys, err = test_helper.SetupSessionForUsers(DB, testUsers[:2])
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test sessions. %v", err)) }

	// Create match request for the rest of users
	_, err = test_helper.SetupMatchRequestForUsers(DB, testUsers[2:])
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Match Request. %v", err)) }

	os.Exit(m.Run())
}