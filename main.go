package main

import (
	"wellnus/backend/references"
	"wellnus/backend/handlers/user"
	"wellnus/backend/handlers/session"
	"wellnus/backend/handlers/group"
	"wellnus/backend/handlers/join"
	
	"fmt"
	"log"
	"github.com/gin-gonic/gin"

	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
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

func main() {
	db := ConnectDB()
	router := gin.Default()

	router.GET("/user", user.GetAllUsersHandler(db))
	router.POST("/user", user.AddUserHandler(db))
	router.GET("/user/:id", user.GetUserHandler(db))
	router.PATCH("/user/:id", user.UpdateUserHandler(db))
	router.DELETE("/user/:id", user.DeleteUserHandler(db))

	router.POST("/session", session.LoginHandler(db))
	router.DELETE("/session", session.LogoutHandler(db))

	router.GET("/group", group.GetAllGroupsHandler(db))
	router.POST("/group", group.AddGroupHandler(db))
	router.GET("/group/:id", group.GetGroupHandler(db))
	router.PATCH("/group/:id", group.UpdateGroupHandler(db))
	router.DELETE("/group/:id", group.LeaveGroupHandler(db))

	router.GET("/join", join.GetAllJoinRequestsHandler(db))
	router.POST("/join", join.AddJoinRequestHandler(db))
	router.GET("/join/:id", join.GetJoinRequestHandler(db))
	router.PATCH("/join/:id", join.RespondJoinRequestHandler(db))
	router.DELETE("/join/:id", join.DeleteJoinRequestHandler(db))

	fmt.Printf("Starting backend server at '%s' \n", references.BACKEND_URL)
	router.Run(references.BACKEND_URL)
}

