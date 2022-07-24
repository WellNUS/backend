package provider

import (
	"wellnus/backend/config"
	"wellnus/backend/db"
	"wellnus/backend/db/model"
	"wellnus/backend/router/provider"
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
type Event = model.Event
type ProviderSetting = model.ProviderSetting

var (
	DB *sql.DB 
	Router *gin.Engine
	NotFoundErrorMessage 		string = http_error.NotFoundError.Error()
	UnauthorizedErrorMessage	string = http_error.UnauthorizedError.Error()
)

// [member, volunteer, counsellor]
var testUsers []User
// [X, ["OffMyChest", "SelfHarm"], ["Anxiety", "SelfHarm"], ["Anxiety", "OffMyChest"]]
var testProviderSettings []ProviderSetting
var testEvents []Event
var sessionKeys []string

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/provider", provider.GetAllProvidersHandler(DB))
	router.GET("/provider/:id", provider.GetProviderWithEventsHandler(DB))
	router.POST("/provider", provider.AddUpdateProviderSettingOfUserHandler(DB))
	router.DELETE("/provider", provider.DeleteProviderSettingOfUserHandler(DB))

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

	testProviderSettings, err = test_helper.SetupProviderSettingForUsers(DB, testUsers[2:])
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test provider setting. %v", err)) }
	// Filling 0, 1 index with empty first
	testProviderSettings = append([]ProviderSetting{ProviderSetting{}, ProviderSetting{}}, testProviderSettings...)

	testEvents, err = test_helper.SetupEventForUsers(DB, testUsers)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test events. %v", err)) }

	os.Exit(m.Run())
}