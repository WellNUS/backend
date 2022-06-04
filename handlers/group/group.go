package group

import (
	"wellnus/backend/config"
	"wellnus/backend/db/model"
	"wellnus/backend/db/query"
	"wellnus/backend/handlers/misc"
	
	"github.com/gin-gonic/gin"
	"database/sql"
)

type Group = model.Group
type GroupWithUsers = model.GroupWithUsers
type User = model.User

// Main functions
func GetAllGroupsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		userID, _ := misc.GetIDCookie(c)
		groups, err := query.GetAllGroups(db, userID)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), groups)
	}
}

func GetGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		groupIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		groupWithUsers, err := query.GetGroupWithUsers(db, groupIDParam)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), groupWithUsers)
	}
}

func AddGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

		newGroup, err := misc.GetGroupFromContext(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}

		newGroup.OwnerID, err = misc.GetIDCookie(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}

		groupWithUsers, err := query.AddGroup(db, newGroup) // Can throw a fatal error
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), groupWithUsers)
	}
}

func UpdateGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		groupIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, err := misc.GetIDCookie(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		updatedGroup, err := misc.GetGroupFromContext(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		updatedGroup, err = query.UpdateGroup(db, updatedGroup, groupIDParam, userIDCookie)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), updatedGroup)
	}
}

func LeaveGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		groupIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, err := misc.GetIDCookie(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		groupWithUsers, err := query.LeaveGroup(db, groupIDParam, userIDCookie)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), groupWithUsers)
	}
}