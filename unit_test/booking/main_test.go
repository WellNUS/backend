package booking

import (
	"wellnus/backend/config"
	"wellnus/backend/db"
	. "wellnus/backend/db/model"
	"wellnus/backend/router/booking"
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

var addedBooking Booking

// [member, volunteer, counsellor]
var testUsers []User
var sessionKeys []string
var testProviderSettings []ProviderSetting
var testBookingsTo2 []Booking
var testBooking0to1 Booking

func setupRouter() *gin.Engine {
	Router := gin.Default()

	Router.GET("/booking", booking.GetAllBookingUsersHandler(DB))
	Router.POST("/booking", booking.AddBookingHandler(DB))
	Router.GET("/booking/:id", booking.GetBookingProviderHandler(DB))
	Router.POST("/booking/:id", booking.RespondBookingHandler(DB))
	Router.PATCH("/booking/:id", booking.UpdateBookingHandler(DB))
	Router.DELETE("/booking/:id", booking.DeleteBookingHandler(DB))

	return Router
}

func TestMain(m *testing.M) {
	config.LoadENV("../../.env")

	DB = db.ConnectDB()
	Router = setupRouter()
	test_helper.ResetDB(DB)
	var err error

	testUsers, err = test_helper.SetupUsers(DB, 3)
	if err != nil {
		log.Fatal(fmt.Sprintf("Something went wrong when creating Test users. %v", err))
	}

	sessionKeys, err = test_helper.SetupSessionForUsers(DB, testUsers)
	if err != nil {
		log.Fatal(fmt.Sprintf("Something went wrong when creating Test sessions. %v", err))
	}

	testProviderSettings, err = test_helper.SetupProviderSettingForUsers(DB, testUsers[2:])
	if err != nil {
		log.Fatal(fmt.Sprintf("Something went wrong when creating Test provider setting. %v", err))
	}
	// Filling 0, 1 index with empty first
	testProviderSettings = append([]ProviderSetting{ProviderSetting{}, ProviderSetting{}}, testProviderSettings...)

	testBookingsTo2, err = test_helper.SetupBookingToUserForUsers(DB, testUsers[:2], testUsers[2])
	if err != nil {
		log.Fatal(fmt.Sprintf("Something went wrong when creating Test Booking. %v", err))
	}

	os.Exit(m.Run())
}
