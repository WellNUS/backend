package group

import (
	"wellnus/backend/config"
	"wellnus/backend/db"
	"wellnus/backend/db/model"
	"wellnus/backend/router/group"
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

type Group = model.Group
type GroupWithUsers = model.GroupWithUsers
type User = model.User

var (
	DB *sql.DB 
	Router *gin.Engine
	NotFoundErrorMessage 		string = http_error.NotFoundError.Error()
	UnauthorizedErrorMessage	string = http_error.UnauthorizedError.Error()
)


var testUsers []User
var sessionKeys []string

var validAddedGroup1 Group = test_helper.GetTestGroup(0)
var validAddedGroup2 Group = Group{
	GroupName: "NewGroupName1",
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/group", group.GetAllGroupsHandler(DB))
	router.POST("/group", group.AddGroupHandler(DB))
	router.DELETE("/group", group.LeaveAllGroupsHandler(DB))
	router.GET("/group/:id", group.GetGroupHandler(DB))
	router.PATCH("/group/:id", group.UpdateGroupHandler(DB))
	router.DELETE("/group/:id", group.LeaveGroupHandler(DB))

	return router
}

func TestMain(m *testing.M) {
	config.LoadENV("../../.env")
	
	DB = db.ConnectTestDB()
	Router = setupRouter()
	test_helper.ResetDB(DB)
	var err error
	
	testUsers, err = test_helper.SetupUsers(DB, 3)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err)) }

	sessionKeys, err = test_helper.SetupSessionForUsers(DB, testUsers)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test sessions. %v", err)) }
	
	os.Exit(m.Run())
}