package match

import (
	"wellnus/backend/db"
	"wellnus/backend/db/model"
	"wellnus/backend/router/match"
	"wellnus/backend/router/misc/http_error"

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
	SessionKey	string
)

var validAddedUser User = User{
	FirstName: "NewFirstName",
	LastName: "NewLastName",
	Gender: "M",
	Faculty: "COMPUTING",
	Email: "NewEmail@u.nus.edu",
	UserRole: "VOLUNTEER",
	Password: "NewPassword",
	PasswordHash: "",
}

var validMatchSetting MatchSetting = MatchSetting{
	FacultyPreference: "MIX",
	Hobbies: []string{"GAMING", "SINGING", "DANCING"},
	MBTI: "INTP",
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/setting", match.GetMatchSettingOfUserHandler(DB))
	router.POST("/setting", match.AddUpdateMatchSettingOfUserHandler(DB))
	router.DELETE("/setting", match.DeleteMatchSettingOfUserHandler(DB))

	router.GET("/match", match.GetLoadedMatchRequestOfUserHandler(DB))
	router.POST("/match", match.AddMatchRequestHandler(DB))
	router.DELETE("/match", match.DeleteMatchRequestOfUserHandler(DB))

	return router
}

func TestMain(m *testing.M) {
	DB = db.ConnectDB()
	Router = setupRouter()
	
	DB.Exec("DELETE FROM wn_group")
	DB.Exec("DELETE FROM wn_user")

	var err error
	validAddedUser, err = model.AddUser(DB, validAddedUser)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err)) }
	SessionKey, err = model.CreateNewSession(DB, validAddedUser.ID)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test sessions. %v", err)) }

	r := m.Run()

	DB.Exec("DELETE FROM wn_user WHERE id = $1", validAddedUser.ID)
	os.Exit(r)
}