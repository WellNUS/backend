package testing

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/unit_test/test_helper"

	"database/sql"
	"net/http"
	"strconv"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetTestingHomeHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		sID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
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
		userID, _ := http_helper.GetIDParams(c)
		userWithGroups, _ := model.GetUserWithGroups(db, userID)
		c.HTML(http.StatusOK, "user.html", gin.H{ "userWithGroups": userWithGroups })
	}
}

func GetTestingAllGroupsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		groups, _ := model.GetAllGroupsOfUser(db, userID)
		c.HTML(http.StatusOK, "groups.html", gin.H{ "groups": groups })
	}
}

func GetTestingGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		groupID, _ := http_helper.GetIDParams(c)
		groupWithUsers, _ := model.GetGroupWithUsers(db, groupID)
		c.HTML(http.StatusOK, "group.html", gin.H{"groupWithUsers": groupWithUsers})
	}
}

func GetTestingAllJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		if s := c.Query("request"); s == "RECEIVED" {
			loadedJoinRequests, _ := model.GetAllLoadedJoinRequestsReceivedOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"loadedJoinRequests": loadedJoinRequests})
		} else if s == "SENT" {
			loadedJoinRequests, _ := model.GetAllLoadedJoinRequestsSentOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"loadedJoinRequests": loadedJoinRequests})
		} else {
			loadedJoinRequests, _ := model.GetAllLoadedJoinRequestsOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"loadedJoinRequests": loadedJoinRequests})
		}
	}
}

func GetTestingJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		joinRequestID, _ := http_helper.GetIDParams(c)
		loadedJoinRequest, _ := model.GetLoadedJoinRequest(db, joinRequestID)
		c.HTML(http.StatusOK, "join.html", gin.H{"loadedJoinRequest": loadedJoinRequest})
	}
}

func GetTestingChatHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		groupID, _ := http_helper.GetIDParams(c)
		groupWithUsers, _ := model.GetGroupWithUsers(db, groupID)
		c.HTML(http.StatusOK, "chat.html", gin.H{"groupWithUsers": groupWithUsers})
	}
}

func GetTestingMatchHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		matchSetting, _ := model.GetMatchSettingOfUser(db, userID)
		count, _ := model.GetMatchRequestCount(db)
		c.HTML(http.StatusOK, "match.html", gin.H{"matchSetting": matchSetting, "mrCount": count})
	}
}

func SetupUsersWithMatchRequests(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		count, err := strconv.Atoi(c.Query("count"))
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
		users, err := test_helper.SetupUsers(db, count)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
		_, err = test_helper.SetupMatchSettingForUsers(db, users)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
		_, err =test_helper.SetupMatchRequestForUsers(db, users)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
		c.IndentedJSON(http.StatusOK, users)
	}
}