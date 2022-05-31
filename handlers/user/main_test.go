package user

// Test should be performed with some users in the database
// Run all test in order of it being written in each file strictly

import (
	"wellnus/backend/references"
	"wellnus/backend/handlers/httpError"

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
	NotFoundErrorMessage 		string = httpError.NotFoundError.Error()
	UnauthorizedErrorMessage	string = httpError.UnauthorizedError.Error()
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

func connectDB() *sql.DB {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
					references.USER,
					references.PASSWORD, 
					references.HOST,
					references.PORT,
					references.DB_NAME)
	// fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Database Connected!")
	return db
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/user", GetAllUsersHandler(db))
	router.POST("/user", AddUserHandler(db))
	router.GET("/user/:id", GetUserHandler(db))
	router.PATCH("/user/:id", UpdateUserHandler(db))
	router.DELETE("/user/:id", DeleteUserHandler(db))

	fmt.Printf("Starting backend server at '%s' \n", references.BACKEND_URL)
	return router
}

func TestMain(m *testing.M) {
	db = connectDB()
	if _, err := db.Query("DELETE FROM wn_user;"); err != nil {
		log.Fatal(fmt.Sprintf("Unable to clear table in preparation for test. %v", err))
	}
	router = setupRouter()
	os.Exit(m.Run())
}