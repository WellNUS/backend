package user

import (
	"wellnus/backend/config"
	"wellnus/backend/db"
	. "wellnus/backend/db/model"
	"wellnus/backend/router/http_helper/http_error"
	"wellnus/backend/router/user"
	"wellnus/backend/unit_test/test_helper"

	"fmt"
	"log"
	"os"
	"testing"

	"database/sql"

	"github.com/gin-gonic/gin"
)

var (
	DB                       *sql.DB
	Router                   *gin.Engine
	addedUser                User
	NotFoundErrorMessage     string = http_error.NotFoundError.Error()
	UnauthorizedErrorMessage string = http_error.UnauthorizedError.Error()
)

var testUsers []User

var validUser User = test_helper.GetTestUser(0)
var sessionKey string

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/user", user.GetAllUsersHandler(DB))
	router.POST("/user", user.AddUserHandler(DB))
	router.GET("/user/:id", user.GetUserHandler(DB))
	router.PATCH("/user/:id", user.UpdateUserHandler(DB))
	router.DELETE("/user/:id", user.DeleteUserHandler(DB))

	return router
}

func TestMain(m *testing.M) {
	config.LoadENV("../../.env")

	DB = db.ConnectDB()
	Router = SetupRouter()
	test_helper.ResetDB(DB)
	var err error

	testUsers, err = test_helper.SetupUsers(DB, 3)
	if err != nil {
		log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err))
	}

	os.Exit(m.Run())
}
