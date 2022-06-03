package group

import (
	"wellnus/backend/references"
	"wellnus/backend/handlers/httpError"
	
	"strconv"
	"github.com/gin-gonic/gin"
	"database/sql"
)

type Group = references.Group
type GroupWithUsers = references.GroupWithUsers
type User = references.User

// Helper functions
func getIDParams(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil { return 0, httpError.NotFoundError }
	return id, nil
}

func getGroupFromContext(c *gin.Context) (Group, error) {
	var group Group
	if err := c.BindJSON(&group); err != nil {
		return Group{}, err
	}
	return group, nil
}

func getIDCookie(c *gin.Context) (int64, error) {
	strUserID, err := c.Cookie("id")
	if err != nil { return 0, httpError.UnauthorizedError }
	userID, err := strconv.ParseInt(strUserID, 0, 64)
	if err != nil { return 0, err }
	return userID, nil
}

// Main functions
func GetAllGroupsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		userID, _ := getIDCookie(c)
		groups, err := GetAllGroups(db, userID)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), groups)
	}
}

func GetGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		groupIDParam, err := getIDParams(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		groupWithUsers, err := GetGroupWithUsers(db, groupIDParam)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), groupWithUsers)
	}
}

func AddGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

		newGroup, err := getGroupFromContext(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}

		newGroup.OwnerID, err = getIDCookie(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}

		groupWithUsers, err := AddGroup(db, newGroup) // Can throw a fatal error
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), groupWithUsers)
	}
}

func UpdateGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		groupIDParam, err := getIDParams(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, err := getIDCookie(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		updatedGroup, err := getGroupFromContext(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		updatedGroup, err = UpdateGroup(db, updatedGroup, groupIDParam, userIDCookie)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), updatedGroup)
	}
}

func LeaveGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		groupIDParam, err := getIDParams(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, err := getIDCookie(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		groupWithUsers, err := LeaveGroup(db, groupIDParam, userIDCookie)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), groupWithUsers)
	}
}