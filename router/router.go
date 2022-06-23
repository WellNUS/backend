package router

import (
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/user"
	"wellnus/backend/router/session"
	"wellnus/backend/router/group"
	"wellnus/backend/router/join"
	"wellnus/backend/router/match"
	"wellnus/backend/router/chat"
	"wellnus/backend/router/testing" //Can be removed at production
	
	"wellnus/backend/router/ws"
	"database/sql"
	
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB, wsHub *ws.Hub) *gin.Engine {
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
	router.GET("/testing/match", testing.GetTestingMatchHandler(db))
	router.POST("/testing/match", testing.SetupUsersWithMatchRequests(db))


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
	
	router.GET("/join", join.GetAllLoadedJoinRequestsHandler(db))
	router.POST("/join", join.AddJoinRequestHandler(db))
	router.GET("/join/:id", join.GetLoadedJoinRequestHandler(db))
	router.PATCH("/join/:id", join.RespondJoinRequestHandler(db))
	router.DELETE("/join/:id", join.DeleteJoinRequestHandler(db))

	router.GET("/setting", match.GetMatchSettingOfUserHandler(db))
	router.POST("/setting", match.AddUpdateMatchSettingOfUserHandler(db))
	router.DELETE("/setting", match.DeleteMatchSettingOfUserHandler(db))

	router.GET("/match", match.GetMatchRequestCount(db))
	router.POST("/match", match.AddMatchRequestHandler(db))
	router.DELETE("/match", match.DeleteMatchRequestOfUserHandler(db))
	router.GET("/match/:id", match.GetLoadedMatchRequestOfUserHandler(db))

	router.GET("/message/:id", chat.GetMessagesChunkOfGroupHandler(db))
	router.GET("/ws/:id", ws.ConnectToWSHandler(wsHub, db))
	
	router.NoRoute(http_helper.NoRouteHandler)

	return router
}