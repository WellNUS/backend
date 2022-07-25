package counsel

import (
	"wellnus/backend/config"
	"wellnus/backend/db"
	"wellnus/backend/db/model"
	"wellnus/backend/router/counsel"
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

type CounselRequest = model.CounselRequest
type User = model.User

var (
	DB *sql.DB 
	Router *gin.Engine
	NotFoundErrorMessage 		string = http_error.NotFoundError.Error()
	UnauthorizedErrorMessage	string = http_error.UnauthorizedError.Error()
)

// Need to assert this
// testUser0 - MEMBER
// testUser1 - VOLUNTEER
// testUser2 - COUNSELLOR

var testUsers []User
var sessionKeys []string
var testCounselRequests []CounselRequest

func setupRouter() *gin.Engine {
	router := gin.Default()
	
	router.GET("/counsel", counsel.GetAllCounselRequestsHandler(DB))
	router.POST("/counsel", counsel.AddUpdateCounselRequestHandler(DB))
	router.DELETE("/counsel", counsel.DeleteCounselRequestHandler(DB))
	router.GET("/counsel/:id", counsel.GetCounselRequestHandler(DB))
	router.POST("/counsel/:id", counsel.AcceptCounselRequestHandler(DB))

	return router
}

func TestMain(m *testing.M) {
	config.LoadENV("../../.env")
	
	DB = db.ConnectDB()
	Router = setupRouter()
	test_helper.ResetDB(DB)
	var err error

	testUsers, err = test_helper.SetupUsers(DB, 3)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err)) }

	sessionKeys, err = test_helper.SetupSessionForUsers(DB, testUsers)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test sessions. %v", err)) }

	testCounselRequests, err = test_helper.SetupCounselRequestForUsers(DB, testUsers[:2])
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test Counsel Request. %v", err)) }

	os.Exit(m.Run())
}