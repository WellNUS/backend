package user

import (
	"wellnus/backend/router/misc"
	"wellnus/backend/router/misc/http_error"
	"wellnus/backend/router/session"
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
		session.CreateNewSessionCookie(db, c, newUser.ID)
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
		userID, _ := misc.GetUserIDFromSessionCookie(db, c)
		if userID != userIDParam {
			err = http_error.UnauthorizedError
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		_, err = model.LeaveAllGroups(db, userID)
		deletedUser, err := model.DeleteUser(db, userID)
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
		userID, _ := misc.GetUserIDFromSessionCookie(db, c)
		if userID != userIDParam {
			err = http_error.UnauthorizedError
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		updatedUser, err := misc.GetUserFromContext(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		updatedUser, err = model.UpdateUser(db,updatedUser, userID)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), updatedUser)
	}
}
