package group

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/misc"
	"wellnus/backend/router/misc/http_error"

	"github.com/gin-gonic/gin"
	"database/sql"
)

// Main functions
func GetAllGroupsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		misc.SetHeaders(c)

		userID, _ := misc.GetUserIDFromSessionCookie(db, c)
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
		misc.SetHeaders(c)

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
		misc.SetHeaders(c)

		newGroup, err := misc.GetGroupFromContext(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		newGroup.OwnerID, err = misc.GetUserIDFromSessionCookie(db, c)
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
		misc.SetHeaders(c)

		groupIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, err := misc.GetUserIDFromSessionCookie(db, c)
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
		misc.SetHeaders(c)

		groupIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, err := misc.GetUserIDFromSessionCookie(db, c)
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
		misc.SetHeaders(c)

		userIDCookie, err := misc.GetUserIDFromSessionCookie(db, c)
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