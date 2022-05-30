package user

// Test should be performed with some users in the database

import (
	"wellnus/backend/references"

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
)

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
	router = setupRouter()
	os.Exit(m.Run())
}