package user

import (
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"
	"wellnus/backend/router/session"
	"wellnus/backend/db/model"
	
	"github.com/gin-gonic/gin"
	"database/sql"
)

type User = model.User

const (
	ROLE_MEMBER = 0
	ROLE_VOLUNTEER = 1
	ROLE_COUNSELLOR = 2
	ROLE_PROVIDER = 3
	ROLE_ALL = 4
)

func getRoleQuery(c *gin.Context) int {
	if s := c.Query("role"); s == "MEMBER" {
		return ROLE_MEMBER
	} else if s == "VOLUNTEER" {
		return ROLE_VOLUNTEER
	} else if s == "COUNSELLOR" {
		return ROLE_COUNSELLOR
	} else if s == "PROVIDER" {
		return ROLE_PROVIDER
	} else {
		return ROLE_ALL
	}
}

// Main functions
func GetAllUsersHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		role := getRoleQuery(c)
		var users []User
		var err error
		if role == ROLE_MEMBER {
			users, err = model.GetAllUsersConditional(db, "user_role = 'MEMBER'")
		}
		if role == ROLE_VOLUNTEER {
			users, err = model.GetAllUsersConditional(db, "user_role = 'VOLUNTEER'")
		}

		if role == ROLE_COUNSELLOR {
			users, err = model.GetAllUsersConditional(db, "user_role = 'COUNSELLOR'")
		}
		
		if role == ROLE_PROVIDER {
			users, err = model.GetAllUsersConditional(db, "user_role = 'VOLUNTEER' OR user_role = 'COUNSELLOR'" )
		}

		if role == ROLE_ALL {
			users, err = model.GetAllUsers(db)
		}

		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), users)
	}
}

func GetUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userIDParam, err := http_helper.GetIDParams(c)
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
		http_helper.SetHeaders(c)

		newUser, err := http_helper.GetUserFromContext(c)
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
		http_helper.SetHeaders(c)

		userIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
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
		http_helper.SetHeaders(c)

		userIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		if userID != userIDParam {
			err = http_error.UnauthorizedError
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		updatedUser, err := http_helper.GetUserFromContext(c)
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
