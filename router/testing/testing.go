package testing

import (
	"wellnus/backend/config"
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/unit_test/test_helper"

	"database/sql"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
)

func GetTestingHomeHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		sID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		c.HTML(http.StatusOK, "home.html", gin.H{ "userID": sID, "backendURL": config.BACKEND_URL})
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
		c.HTML(http.StatusOK, "user.html", gin.H{ "userWithGroups": userWithGroups, "backendURL": config.BACKEND_URL })
	}
}

func GetTestingAllGroupsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		groups, _ := model.GetAllGroupsOfUser(db, userID)
		c.HTML(http.StatusOK, "groups.html", gin.H{ "groups": groups, "backendURL": config.BACKEND_URL })
	}
}

func GetTestingGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		groupID, _ := http_helper.GetIDParams(c)
		groupWithUsers, _ := model.GetGroupWithUsers(db, groupID)
		c.HTML(http.StatusOK, "group.html", gin.H{"groupWithUsers": groupWithUsers, "backendURL": config.BACKEND_URL})
	}
}

func GetTestingAllJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		if s := c.Query("request"); s == "RECEIVED" {
			loadedJoinRequests, _ := model.GetAllLoadedJoinRequestsReceivedOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"loadedJoinRequests": loadedJoinRequests, "backendURL": config.BACKEND_URL})
		} else if s == "SENT" {
			loadedJoinRequests, _ := model.GetAllLoadedJoinRequestsSentOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"loadedJoinRequests": loadedJoinRequests, "backendURL": config.BACKEND_URL})
		} else {
			loadedJoinRequests, _ := model.GetAllLoadedJoinRequestsOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"loadedJoinRequests": loadedJoinRequests, "backendURL": config.BACKEND_URL})
		}
	}
}

func GetTestingJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		joinRequestID, _ := http_helper.GetIDParams(c)
		loadedJoinRequest, _ := model.GetLoadedJoinRequest(db, joinRequestID)
		c.HTML(http.StatusOK, "join.html", gin.H{"loadedJoinRequest": loadedJoinRequest, "backendURL": config.BACKEND_URL})
	}
}

func GetTestingChatHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		groupID, _ := http_helper.GetIDParams(c)
		groupWithUsers, _ := model.GetGroupWithUsers(db, groupID)
		c.HTML(http.StatusOK, "chat.html", gin.H{"groupWithUsers": groupWithUsers, "backendURL": config.BACKEND_URL})
	}
}

func GetTestingMatchHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		matchSetting, _ := model.GetMatchSettingOfUser(db, userID)
		count, _ := model.GetMatchRequestCount(db)
		c.HTML(http.StatusOK, "match.html", gin.H{"matchSetting": matchSetting, "mrCount": count, "backendURL": config.BACKEND_URL})
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

func GetTestingCounselRequestsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		topics, _ := c.GetQueryArray("topic")
		counselRequest, _ := model.GetCounselRequest(db, userID, userID)
		counselRequests, _ := model.GetAllCounselRequests(db, topics, userID)
		c.HTML(http.StatusOK, "counsel_requests.html", gin.H{"counselRequests": counselRequests, "counselRequest": counselRequest, "backendURL": config.BACKEND_URL})
	}
}

func GetTestingCounselRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userIDParam, _ := http_helper.GetIDParams(c)
		userIDCookie, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		counselRequest, _ := model.GetCounselRequest(db, userIDParam, userIDCookie)
		c.HTML(http.StatusOK, "counsel_request.html", gin.H{"counselRequest": counselRequest, "backendURL": config.BACKEND_URL})
	}
}