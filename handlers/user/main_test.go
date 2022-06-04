package user

import (
	"wellnus/backend/db/model"
	"wellnus/backend/config"
	"wellnus/backend/handlers/misc"

	"testing"
	"os"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"database/sql"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB 
	router *gin.Engine
	addedUser User
	NotFoundErrorMessage 		string = misc.NotFoundError.Error()
	UnauthorizedErrorMessage	string = misc.UnauthorizedError.Error()
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

func equal(user1 User, user2 User) bool {
	return user1.ID == user2.ID &&
		user1.FirstName == user2.FirstName &&
		user1.LastName == user2.LastName &&
		user1.Gender == user2.Gender &&
		user1.Faculty == user2.Faculty &&
		user1.Email == user2.Email &&
		user1.UserRole == user2.UserRole &&
		user1.PasswordHash == user2.PasswordHash
}

func setupDB() *sql.DB {
	db, err := sql.Open("postgres", config.CONNECTION_STRING)
	if err != nil {
		log.Fatal(err)
	}
	db.Query("DELETE FROM wn_group;")
	db.Query("DELETE FROM wn_user;")
	return db
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/user", GetAllUsersHandler(db))
	router.POST("/user", AddUserHandler(db))
	router.GET("/user/:id", GetUserHandler(db))
	router.PATCH("/user/:id", UpdateUserHandler(db))
	router.DELETE("/user/:id", DeleteUserHandler(db))

	fmt.Printf("Starting backend server at '%s' \n", config.BACKEND_URL)
	return router
}

func TestMain(m *testing.M) {
	db = setupDB()
	router = setupRouter()
	os.Exit(m.Run())
}