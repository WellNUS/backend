package match

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
var validMatchSetting MatchSetting = test_helper.GetRandomTestMatchSetting()

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/setting", match.GetMatchSettingOfUserHandler(DB))
	router.POST("/setting", match.AddUpdateMatchSettingOfUserHandler(DB))
	router.DELETE("/setting", match.DeleteMatchSettingOfUserHandler(DB))

	router.GET("/match", match.GetMatchRequestCount(DB))
	router.POST("/match", match.AddMatchRequestHandler(DB))
	router.DELETE("/match", match.DeleteMatchRequestOfUserHandler(DB))
	router.GET("/match/:id", match.GetLoadedMatchRequestOfUserHandler(DB))

	return router
}

func TestMain(m *testing.M) {
	config.LoadENV("../../.env")
	
	DB = db.ConnectDB()
	Router = setupRouter()
	test_helper.ResetDB(DB)
	var err error

	testUsers, err = test_helper.SetupUsers(DB, 1)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err)) }

	sessionKeys, err = test_helper.SetupSessionForUsers(DB, testUsers)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test sessions. %v", err)) }	

	os.Exit(m.Run())
}