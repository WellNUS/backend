package group

import (
	"wellnus/backend/config"
	"wellnus/backend/db/model"
	"wellnus/backend/router/misc"
	"wellnus/backend/router/misc/http_error"

	"github.com/gin-gonic/gin"
	"database/sql"
)

// Main functions
func GetAllGroupsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")
		userID, _ := misc.GetIDCookie(c)
		groups, err := model.GetAllGroupsOfUser(db, userID)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), groups)
	}
}

func GetGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")
		groupIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		groupWithUsers, err := model.GetGroupWithUsers(db, groupIDParam)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), groupWithUsers)
	}
}

func AddGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")

		newGroup, err := misc.GetGroupFromContext(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		newGroup.OwnerID, err = misc.GetIDCookie(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		groupWithUsers, err := model.AddGroup(db, newGroup) // Can throw a fatal error
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), groupWithUsers)
	}
}

func UpdateGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")
		groupIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, err := misc.GetIDCookie(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		updatedGroup, err := misc.GetGroupFromContext(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		updatedGroup, err = model.UpdateGroup(db, updatedGroup, groupIDParam, userIDCookie)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), updatedGroup)
	}
}

func LeaveGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")

		groupIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, err := misc.GetIDCookie(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		groupWithUsers, err := model.LeaveGroup(db, groupIDParam, userIDCookie)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), groupWithUsers)
	}
}

func LeaveAllGroupsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")
		userIDCookie, err := misc.GetIDCookie(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		groupsWithUsers, err := model.LeaveAllGroups(db, userIDCookie)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), groupsWithUsers)
	}
}