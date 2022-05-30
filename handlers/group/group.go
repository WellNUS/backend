package group

import (
	"wellnus/backend/references"
	"wellnus/backend/handlers/httpError"
	
	"fmt"
	"errors"
	"strconv"
	"github.com/gin-gonic/gin"
	"database/sql"
)

type Group = references.Group
type GroupWithUsers = references.GroupWithUsers
type User = references.User

// Helper functions
func GetIDParams(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 0, 64)
}

func GetGroupFromContext(c *gin.Context) (Group, error) {
	var group Group
	if err := c.BindJSON(&group); err != nil {
		return Group{}, err
	}
	return group, nil
}

// func GetUserFromContext(c *gin.Context)

func GetIDCookie(c *gin.Context) (int64, error) {
	strUserID, err := c.Cookie("id")
	if err != nil { return 0, err }
	userID, err := strconv.ParseInt(strUserID, 0, 64)
	if err != nil { return 0, err }
	return userID, nil
}

// Main functions
func GetAllGroupsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		groups, err := GetAllGroups(db)
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
		id, err := GetIDParams(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		groupWithUsers, err := GetGroup(db, id)
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

		newGroup, err := GetGroupFromContext(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}

		newGroup.OwnerID, err = GetIDCookie(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}

		newGroup, err = AddGroup(db, newGroup) // Can throw a fatal error
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}

		// newGroup has been added properly. Look for users in group
		users, err := GetUsersInGroup(db, newGroup.ID)
		if err != nil {
			err = errors.New(fmt.Sprintf("New group has been added into database, but failed to get users in new group. %v", err))
			c.IndentedJSON(httpError.GetStatusCode(err), err)
		}
		c.IndentedJSON(httpError.GetStatusCode(err), GroupWithUsers{ Group: newGroup, Users: users })
	}
}

func AddUserToGroupHandler(db *sql.DB) func(*gin.Context) { // Might be deprecated when request system is put up
	return func(c *gin.Context) {
		c.IndentedJSON(200, "hello")
	}
}