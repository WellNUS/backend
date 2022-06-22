package join

import (
	"wellnus/backend/db"
	"wellnus/backend/db/model"
	"wellnus/backend/router/join"
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
type Group = model.Group
type JoinRequest = model.JoinRequest
type LoadedJoinRequest = model.LoadedJoinRequest
type JoinRequestRespond = model.JoinRequestRespond

var (
	DB *sql.DB 
	Router *gin.Engine
	NotFoundErrorMessage 		string = http_error.NotFoundError.Error()
	UnauthorizedErrorMessage	string = http_error.UnauthorizedError.Error()
	SessionKey1		string
	SessionKey2		string
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

var validAddedGroup Group = Group{
	GroupName: "NewGroupName",
	GroupDescription: "NewGroupDescription",
	Category: "SUPPORT",
}

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
	DB = db.ConnectDB()
	Router = setupRouter()
	
	DB.Exec("DELETE FROM wn_group")
	DB.Exec("DELETE FROM wn_user")

	var err error
	validAddedUser1, err = model.AddUser(DB, validAddedUser1)
	validAddedUser2, err = model.AddUser(DB, validAddedUser2)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err)) }
	validAddedGroup.OwnerID = validAddedUser1.ID	//Setting user1 as owner
	validAddedGroupWithUser, err := model.AddGroup(DB, validAddedGroup)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test group. %v", err)) }
	validAddedGroup	= validAddedGroupWithUser.Group

	SessionKey1, err = model.CreateNewSession(DB, validAddedUser1.ID)
	SessionKey2, err = model.CreateNewSession(DB, validAddedUser2.ID)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test sessions. %v", err)) }


	r := m.Run()

	DB.Exec("DELETE FROM wn_group WHERE id = $1", validAddedGroup.ID)
	DB.Exec("DELETE FROM wn_user WHERE id = $1 OR id = $2", validAddedUser1.ID, validAddedUser2.ID)
	os.Exit(r)
}