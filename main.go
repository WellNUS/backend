package main

import (
	"wellnus/backend/config"
	"wellnus/backend/handlers/user"
	"wellnus/backend/handlers/session"
	"wellnus/backend/handlers/group"
	"wellnus/backend/handlers/join"
	"wellnus/backend/handlers/testing" //Can be removed at production

	"wellnus/backend/handlers/ws"
	
	"fmt"
	"log"
	"github.com/gin-gonic/gin"

	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	// fmt.Println(connStr)
	db, err := sql.Open("postgres", config.CONNECTION_STRING)
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
	wsHub := ws.NewHub(db)
	go wsHub.Run()
	
	router := gin.Default()

	// Remove this on production
	router.LoadHTMLGlob("templates/**/*")
	router.GET("/testing", testing.GetTestingHome(db))
	router.GET("/testing/user", testing.GetTestingAllUsers(db))
	router.GET("/testing/user/:id", testing.GetTestingUser(db))
	router.GET("/testing/group", testing.GetTestingAllGroups(db))
	router.GET("/testing/group/:id", testing.GetTestingGroup(db))
	router.GET("/testing/join", testing.GetTestingAllJoinRequest(db))
	router.GET("/testing/join/:id", testing.GetTestingJoinRequest(db))

	router.GET("/user", user.GetAllUsersHandler(db))
	router.POST("/user", user.AddUserHandler(db))
	router.GET("/user/:id", user.GetUserHandler(db))
	router.PATCH("/user/:id", user.UpdateUserHandler(db))
	router.DELETE("/user/:id", user.DeleteUserHandler(db))

	router.POST("/session", session.LoginHandler(db))
	router.DELETE("/session", session.LogoutHandler(db))

	router.GET("/group", group.GetAllGroupsHandler(db))
	router.POST("/group", group.AddGroupHandler(db))
	router.DELETE("/group", group.LeaveAllGroupsHandler(db))
	router.GET("/group/:id", group.GetGroupHandler(db))
	router.PATCH("/group/:id", group.UpdateGroupHandler(db))
	router.DELETE("/group/:id", group.LeaveGroupHandler(db))

	router.GET("/join", join.GetAllJoinRequestsHandler(db))
	router.POST("/join", join.AddJoinRequestHandler(db))
	router.GET("/join/:id", join.GetLoadedJoinRequestHandler(db))
	router.PATCH("/join/:id", join.RespondJoinRequestHandler(db))
	router.DELETE("/join/:id", join.DeleteJoinRequestHandler(db))

	router.GET("/ws", ws.ConnectToWSHandler(wsHub))

	fmt.Printf("Starting backend server at '%s' \n", config.BACKEND_URL)
	router.Run(config.BACKEND_URL)
}

