package testing

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/misc"

	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTestingHomeHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		sID, _ := misc.GetUserIDFromSessionCookie(db, c)
		c.HTML(http.StatusOK, "home.html", gin.H{ "userID": sID })
	}
}

func GetTestingAllUsersHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		users, _ := model.GetAllUsers(db)
		c.HTML(http.StatusOK, "users.html", gin.H{ "users": users })
	}
}

func GetTestingUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := misc.GetIDParams(c)
		userWithGroups, _ := model.GetUserWithGroups(db, userID)
		c.HTML(http.StatusOK, "user.html", gin.H{ "userWithGroups": userWithGroups })
	}
}

func GetTestingAllGroupsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := misc.GetUserIDFromSessionCookie(db, c)
		groups, _ := model.GetAllGroupsOfUser(db, userID)
		c.HTML(http.StatusOK, "groups.html", gin.H{ "groups": groups })
	}
}

func GetTestingGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		groupID, _ := misc.GetIDParams(c)
		groupWithUsers, _ := model.GetGroupWithUsers(db, groupID)
		c.HTML(http.StatusOK, "group.html", gin.H{"groupWithUsers": groupWithUsers})
	}
}

func GetTestingAllJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := misc.GetUserIDFromSessionCookie(db, c)
		if s := c.Query("request"); s == "RECEIVED" {
			joinRequests, _ := model.GetAllJoinRequestsReceivedOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"joinRequests": joinRequests})
		} else if s == "SENT" {
			joinRequests, _ := model.GetAllJoinRequestsSentOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"joinRequests": joinRequests})
		} else {
			joinRequests, _ := model.GetAllJoinRequestsOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"joinRequests": joinRequests})
		}
	}
}

func GetTestingJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		joinRequestID, _ := misc.GetIDParams(c)
		loadedJoinRequest, _ := model.GetLoadedJoinRequest(db, joinRequestID)
		c.HTML(http.StatusOK, "join.html", gin.H{"loadedJoinRequest": loadedJoinRequest})
	}
}

func GetTestingChatHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		groupID, _ := misc.GetIDParams(c)
		groupWithUsers, _ := model.GetGroupWithUsers(db, groupID)
		c.HTML(http.StatusOK, "chat.html", gin.H{"groupWithUsers": groupWithUsers})
	}
}