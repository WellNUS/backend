package testing

import (
	"wellnus/backend/db/query"
	"wellnus/backend/handlers/misc"

	"database/sql"
	//"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTestingHome(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		sID, _ := c.Cookie("id")
		c.HTML(http.StatusOK, "home.html", gin.H{ "userID": sID })
	}
}

func GetTestingAllUsers(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		users, _ := query.GetAllUsers(db)
		c.HTML(http.StatusOK, "users.html", gin.H{ "users": users })
	}
}

func GetTestingUser(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := misc.GetIDParams(c)
		userWithGroups, _ := query.GetUserWithGroups(db, userID)
		c.HTML(http.StatusOK, "user.html", gin.H{ "userWithGroups": userWithGroups })
	}
}

func GetTestingAllGroups(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := misc.GetIDCookie(c)
		groups, _ := query.GetAllGroupsOfUser(db, userID)
		c.HTML(http.StatusOK, "groups.html", gin.H{ "groups": groups })
	}
}

func GetTestingGroup(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		groupID, _ := misc.GetIDParams(c)
		groupWithUsers, _ := query.GetGroupWithUsers(db, groupID)
		c.HTML(http.StatusOK, "group.html", gin.H{"groupWithUsers": groupWithUsers})
	}
}

func GetTestingAllJoinRequest(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := misc.GetIDCookie(c)
		if s := c.Query("request"); s == "RECEIVED" {
			joinRequests, _ := query.GetAllJoinRequestsReceivedOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"joinRequests": joinRequests})
		} else if s == "SENT" {
			joinRequests, _ := query.GetAllJoinRequestsSentOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"joinRequests": joinRequests})
		} else {
			joinRequests, _ := query.GetAllJoinRequestsOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"joinRequests": joinRequests})
		}
	}
}

func GetTestingJoinRequest(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		joinRequestID, _ := misc.GetIDParams(c)
		loadedJoinRequest, _ := query.GetLoadedJoinRequest(db, joinRequestID)
		c.HTML(http.StatusOK, "join.html", gin.H{"loadedJoinRequest": loadedJoinRequest})
	}
}