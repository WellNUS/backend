package event

import (
	"wellnus/backend/config"
	"wellnus/backend/db"
	. "wellnus/backend/db/model"
	"wellnus/backend/router/event"
	"wellnus/backend/router/http_helper/http_error"
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
	DB                       *sql.DB
	Router                   *gin.Engine
	NotFoundErrorMessage     string = http_error.NotFoundError.Error()
	UnauthorizedErrorMessage string = http_error.UnauthorizedError.Error()
)

var testUsers []User
var sessionKeys []string

// testEvents[0] -> public, testEvents[1] -> private
var testEvents []Event

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/event", event.GetAllEventsHandler(DB))
	router.POST("/event", event.AddEventHandler(DB))
	router.DELETE("/event", event.LeaveDeleteAllEventsHandler(DB))
	router.GET("/event/:id", event.GetEventHandler(DB))
	router.POST("/event/:id", event.AddUserToEventHandler(DB))
	router.PATCH("/event/:id", event.UpdateEventHandler(DB))
	router.DELETE("event/:id", event.LeaveDeleteEventHandler(DB))
	router.POST("/event/:id/start", event.CreateGroupDeleteEventHandler(DB))

	return router
}

func TestMain(m *testing.M) {
	config.LoadENV("../../.env")

	DB = db.ConnectDB()
	Router = setupRouter()
	test_helper.ResetDB(DB)
	var err error

	testUsers, err = test_helper.SetupUsers(DB, 2)
	if err != nil {
		log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err))
	}

	sessionKeys, err = test_helper.SetupSessionForUsers(DB, testUsers)
	if err != nil {
		log.Fatal(fmt.Sprintf("Something went wrong when creating Test sessions. %v", err))
	}

	testEvents, err = test_helper.SetupEventForUsers(DB, testUsers[:1])
	if err != nil {
		log.Fatal(fmt.Sprintf("Something went wrong when creating Test Events. %v", err))
	}

	os.Exit(m.Run())
}
