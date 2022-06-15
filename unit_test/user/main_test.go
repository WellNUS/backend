package user

import (
	"wellnus/backend/db"
	"wellnus/backend/db/model"
	"wellnus/backend/router/user"
	"wellnus/backend/router/misc/http_error"
	
	"testing"
	"os"

	"database/sql"
	"github.com/gin-gonic/gin"
)

type User = model.User
type UserWithGroups = model.UserWithGroups

var (
	DB *sql.DB 
	Router *gin.Engine
	addedUser User
	NotFoundErrorMessage 		string = http_error.NotFoundError.Error()
	UnauthorizedErrorMessage	string = http_error.UnauthorizedError.Error()
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
	DB = db.ConnectDB()
	Router = SetupRouter()

	DB.Exec("DELETE FROM wn_group")
	DB.Exec("DELETE FROM wn_user")

	os.Exit(m.Run())
}