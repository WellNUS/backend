package join

import (
	"wellnus/backend/config"
	"wellnus/backend/db"
	"wellnus/backend/db/model"
	"wellnus/backend/router/join"
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
type Group = model.Group
type JoinRequest = model.JoinRequest
type LoadedJoinRequest = model.LoadedJoinRequest
type JoinRequestRespond = model.JoinRequestRespond

var (
	DB *sql.DB 
	Router *gin.Engine
	NotFoundErrorMessage 		string = http_error.NotFoundError.Error()
	UnauthorizedErrorMessage	string = http_error.UnauthorizedError.Error()
)

var addedJoinRequest JoinRequest
var testUsers []User
var sessionKeys []string
var testGroups []Group

func setupRouter() *gin.Engine {
	Router := gin.Default()

	Router.GET("/join", join.GetAllLoadedJoinRequestsHandler(DB))
	Router.POST("/join", join.AddJoinRequestHandler(DB))
	Router.GET("/join/:id", join.GetLoadedJoinRequestHandler(DB))
	Router.PATCH("/join/:id", join.RespondJoinRequestHandler(DB))
	Router.DELETE("/join/:id", join.DeleteJoinRequestHandler(DB))

	return Router
}

func TestMain(m *testing.M) {
	config.LoadENV("../../.env")
	
	DB = db.ConnectTestDB()
	Router = setupRouter()
	test_helper.ResetDB(DB)
	var err error

	testUsers, err = test_helper.SetupUsers(DB, 2)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test users. %v", err)) }

	testGroups, err = test_helper.SetupGroupsForUsers(DB, testUsers[:1])
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test group. %v", err)) }

	sessionKeys, err = test_helper.SetupSessionForUsers(DB, testUsers)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test sessions. %v", err)) }

	os.Exit(m.Run())
}