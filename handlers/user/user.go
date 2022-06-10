package user

import (
	"wellnus/backend/config"
	"wellnus/backend/handlers/misc"
	"wellnus/backend/db/model"
	"wellnus/backend/db/query"
	
	"github.com/gin-gonic/gin"
	"database/sql"
)

type User = model.User
type UserWithGroups = model.UserWithGroups

// Main functions
func GetAllUsersHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		users, err := query.GetAllUsers(db)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), users)
	}
}

func GetUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		userIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		userWithGroups, err := query.GetUserWithGroups(db, userIDParam)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), userWithGroups)
	}
}

func AddUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		newUser, err := misc.GetUserFromContext(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		newUser, err = query.AddUser(db, newUser)
		if err != nil { 
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		misc.SetIDCookie(c, newUser.ID)
		c.IndentedJSON(misc.GetStatusCode(err), newUser)
	}
}

func DeleteUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		userIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, _ := misc.GetIDCookie(c)
		if userIDCookie != userIDParam {
			err = misc.UnauthorizedError
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		_, err = query.LeaveAllGroups(db, userIDCookie)
		deletedUser, err := query.DeleteUser(db, userIDCookie)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), deletedUser)
	}
}

func UpdateUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		userIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, _ := misc.GetIDCookie(c)
		if userIDCookie != userIDParam {
			err = misc.UnauthorizedError
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		updatedUser, err := misc.GetUserFromContext(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		updatedUser, err = query.UpdateUser(db, updatedUser, userIDCookie)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), updatedUser)
	}
}
