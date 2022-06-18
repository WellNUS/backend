package session

import (
	"wellnus/backend/db"
	"wellnus/backend/db/model"
	"wellnus/backend/router/session"

	"testing"
	"os"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"database/sql"
	_ "github.com/lib/pq"
)

type User = model.User
type SessionResponse = model.SessionResponse

var (
	DB *sql.DB
	Router *gin.Engine
	sessionKey1 string
	sessionKey2 string
) 

var validUser User = User{
	FirstName: "NewFirstName",
	LastName: "NewLastName",
	Gender: "M",
	Faculty: "COMPUTING",
	Email: "NewEmail@u.nus.edu",
	UserRole: "VOLUNTEER",
	Password: "NewPassword",
	PasswordHash: "",
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	
	router.POST("/session", session.LoginHandler(DB))
	router.DELETE("/session", session.LogoutHandler(DB))

	return router
}

func TestMain(m *testing.M) {
	DB = db.ConnectDB()
	Router = setupRouter()

	DB.Exec("DELETE FROM wn_group")
	DB.Exec("DELETE FROM wn_user")

	user, err := model.AddUser(DB, validUser)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err)) }

	r := m.Run()

	DB.Exec("DELETE FROM wn_user WHERE id = $1", user.ID)
	os.Exit(r)
}