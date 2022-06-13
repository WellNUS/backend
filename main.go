package main

import (
	"wellnus/backend/config"
	"wellnus/backend/handlers/user"
	"wellnus/backend/handlers/session"
	"wellnus/backend/handlers/group"
	"wellnus/backend/handlers/join"
	"wellnus/backend/handlers/message"
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
	router.GET("/testing", testing.GetTestingHomeHandler(db))
	router.GET("/testing/user", testing.GetTestingAllUsersHandler(db))
	router.GET("/testing/user/:id", testing.GetTestingUserHandler(db))
	router.GET("/testing/group", testing.GetTestingAllGroupsHandler(db))
	router.GET("/testing/group/:id", testing.GetTestingGroupHandler(db))
	router.GET("/testing/group/:id/chat", testing.GetTestingChatHandler(db))
	router.GET("/testing/join", testing.GetTestingAllJoinRequestHandler(db))
	router.GET("/testing/join/:id", testing.GetTestingJoinRequestHandler(db))
	

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

	router.GET("/message/:id", message.GetLoadedMessagesPacketOfGroupHandler(db))

	router.GET("/ws/:id", ws.ConnectToWSHandler(wsHub, db))

	fmt.Printf("Starting backend server at '%s' \n", config.BACKEND_URL)
	router.Run(config.BACKEND_URL)
}

