package user

import (
	"wellnus/backend/router/misc"
	"wellnus/backend/router/misc/http_error"
	"wellnus/backend/db/model"
	
	"github.com/gin-gonic/gin"
	"database/sql"
)

// Main functions
func GetAllUsersHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		misc.SetHeaders(c)

		users, err := model.GetAllUsers(db)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), users)
	}
}

func GetUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		misc.SetHeaders(c)

		userIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userWithGroups, err := model.GetUserWithGroups(db, userIDParam)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), userWithGroups)
	}
}

func AddUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		misc.SetHeaders(c)

		newUser, err := misc.GetUserFromContext(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		newUser, err = model.AddUser(db, newUser)
		if err != nil { 
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		misc.SetIDCookie(c, newUser.ID)
		c.IndentedJSON(http_error.GetStatusCode(err), newUser)
	}
}

func DeleteUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		misc.SetHeaders(c)

		userIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, _ := misc.GetIDCookie(c)
		if userIDCookie != userIDParam {
			err = http_error.UnauthorizedError
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		_, err = model.LeaveAllGroups(db, userIDCookie)
		deletedUser, err := model.DeleteUser(db, userIDCookie)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), deletedUser)
	}
}

func UpdateUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		misc.SetHeaders(c)

		userIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, _ := misc.GetIDCookie(c)
		if userIDCookie != userIDParam {
			err = http_error.UnauthorizedError
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		updatedUser, err := misc.GetUserFromContext(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		updatedUser, err = model.UpdateUser(db,updatedUser, userIDCookie)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), updatedUser)
	}
}
