package user

import (
	"wellnus/backend/references"
	"wellnus/backend/handlers/httpError"
	
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"
	"database/sql"
)

type User = references.User

// Helper functions
func getIDParams(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 0, 64)
}

func getIDCookie(c *gin.Context) (int64, error) {
	strUserID, err := c.Cookie("id")
	if err != nil { return 0, err }
	userID, err := strconv.ParseInt(strUserID, 0, 64)
	if err != nil { return 0, err }
	return userID, nil
}

func setIDCookie(c *gin.Context, id int64) {
	c.SetCookie("id", fmt.Sprintf("%d", id), 1209600, "/", references.DOMAIN, false, true)
}

func getUserFromContext(c *gin.Context) (User, error) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		return User{}, nil
	}
	return user, nil
}

// Main functions

func GetAllUsersHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		users, err := GetAllUsers(db)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), users)
	}
}

func GetUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		id, err := getIDParams(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		user, err := GetUser(db, id)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), user)
	}
}

func AddUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		newUser, err := getUserFromContext(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		newUser, err = AddUser(db, newUser)
		if err != nil { 
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		setIDCookie(c, newUser.ID)
		c.IndentedJSON(httpError.GetStatusCode(err), newUser)
	}
}

func DeleteUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		id, err := getIDParams(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		userID, _ := getIDCookie(c)
		if userID != id {
			err = httpError.UnauthorizedError
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		deletedUser, err := DeleteUser(db, id)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), deletedUser)
	}
}

func UpdateUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		id, err := getIDParams(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		userID, _ := getIDCookie(c)
		if userID != id {
			err = httpError.UnauthorizedError
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		updatedUser, err := getUserFromContext(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		updatedUser, err = UpdateUser(db, updatedUser, id)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), updatedUser)
	}
}
