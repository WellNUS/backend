package group

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"

	"github.com/gin-gonic/gin"
	"database/sql"
)

// Main functions
func GetAllGroupsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		groups, err := model.GetAllGroupsOfUser(db, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), groups)
	}
}

func GetGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		groupIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		groupWithUsers, err := model.GetGroupWithUsers(db, groupIDParam)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), groupWithUsers)
	}
}

func AddGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		newGroup, err := http_helper.GetGroupFromContext(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		newGroup.Category = "CUSTOM"

		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		groupWithUsers, err := model.AddGroupWithUserIDs(db, newGroup, []int64{userID}) // Can throw a fatal error
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), groupWithUsers)
	}
}

func UpdateGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		groupIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		updatedGroup, err := http_helper.GetGroupFromContext(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		updatedGroup, err = model.UpdateGroup(db, updatedGroup, groupIDParam, userIDCookie)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), updatedGroup)
	}
}

func LeaveGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		groupIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		groupWithUsers, err := model.LeaveGroup(db, groupIDParam, userIDCookie)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), groupWithUsers)
	}
}

func LeaveAllGroupsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userIDCookie, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		groupsWithUsers, err := model.LeaveAllGroups(db, userIDCookie)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), groupsWithUsers)
	}
}