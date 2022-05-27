package main

import (
	"wellnus/backend/references"
	"wellnus/backend/handlers/users"
	"wellnus/backend/handlers/session"
	
	"fmt"
	"log"
	"github.com/gin-gonic/gin"

	"database/sql"
	_ "github.com/lib/pq"
)

func connectDB() *sql.DB {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
					references.USER,
					references.PASSWORD, 
					references.HOST,
					references.PORT,
					references.DB_NAME)
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

func main() {
	db := connectDB()
	router := gin.Default()

	router.GET("/users", users.GetAllUsersHandler(db))
	router.POST("/users", users.AddUserHandler(db))
	router.GET("/users/:id", users.GetUserHandler(db))
	router.PATCH("/users/:id", users.UpdateUserHandler(db))
	router.DELETE("/users/:id", users.DeleteUserHandler(db))

	router.POST("/session", session.LoginHandler(db))
	router.DELETE("/session", session.LogoutHandler(db))

	fmt.Printf("Starting backend server at '%s' \n", references.BACKEND_URL)
	router.Run(references.BACKEND_URL)
}

