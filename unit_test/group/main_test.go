package group

import (
	"wellnus/backend/db"
	"wellnus/backend/db/model"
	"wellnus/backend/router/group"
	"wellnus/backend/router/misc/http_error"

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
	SessionKey1	string
	SessionKey2 string
)

var validAddedUser1 User = User{
	FirstName: "NewFirstName",
	LastName: "NewLastName",
	Gender: "M",
	Faculty: "COMPUTING",
	Email: "NewEmail@u.nus.edu",
	UserRole: "VOLUNTEER",
	Password: "NewPassword",
	PasswordHash: "",
}

var validAddedUser2 User = User{
	FirstName: "NewFirstName1",
	LastName: "NewLastName1",
	Gender: "M",
	Faculty: "COMPUTING",
	Email: "NewEmail1@u.nus.edu",
	UserRole: "VOLUNTEER",
	Password: "NewPassword",
	PasswordHash: "",
}

var validAddedGroup1 Group = Group{
	GroupName: "NewGroupName",
	GroupDescription: "NewGroupDescription",
	Category: "SUPPORT",
}

var validAddedGroup2 Group = Group{
	GroupName: "NewGroupName1",
	Category: "SUPPORT",
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
	DB = db.ConnectDB()
	Router = setupRouter()
	
	DB.Exec("DELETE FROM wn_group")
	DB.Exec("DELETE FROM wn_user")

	var err error
	validAddedUser1, err = model.AddUser(DB, validAddedUser1)
	validAddedUser2, err = model.AddUser(DB, validAddedUser2)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err)) }
	SessionKey1, err = model.CreateNewSession(DB, validAddedUser1.ID)
	SessionKey2, err = model.CreateNewSession(DB, validAddedUser2.ID)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test sessions. %v", err)) }

	r := m.Run()

	DB.Exec("DELETE FROM wn_user WHERE id = $1 OR id = $2", validAddedUser1.ID, validAddedUser2.ID)
	os.Exit(r)
}