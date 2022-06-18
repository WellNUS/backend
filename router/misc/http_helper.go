package misc

import (
	"wellnus/backend/db/model"
	"wellnus/backend/config"
	"wellnus/backend/router/misc/http_error"
	"strconv"

	"database/sql"
	"github.com/gin-gonic/gin"
)

type User = model.User
type Group = model.Group
type JoinRequest = model.JoinRequest
type JoinRequestRespond = model.JoinRequestRespond

func SetHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
	c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
	c.Header("Access-Control-Allow-Credentials", "true")
}

func GetIDParams(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil { return 0, http_error.NotFoundError }
	return id, nil
}

func GetUserIDFromSessionCookie(db *sql.DB, c *gin.Context) (int64, error) {
	sessionKey, err := c.Cookie("session_key")
	if err != nil { return 0, http_error.UnauthorizedError }
	userID, err := model.GetUserIDFromSessionKey(db, sessionKey)
	if err != nil { return 0, err }
	return userID, nil
}


func GetUserFromContext(c *gin.Context) (User, error) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		return User{}, nil
	}
	return user, nil
}

func GetJoinRequestFromContext(c *gin.Context) (JoinRequest, error) {
	var joinRequest JoinRequest
	if err := c.BindJSON(&joinRequest); err != nil {
		return JoinRequest{}, err
	}
	return joinRequest, nil
}

func GetGroupFromContext(c *gin.Context) (Group, error) {
	var group Group
	if err := c.BindJSON(&group); err != nil {
		return Group{}, err
	}
	return group, nil
}

func GetJoinRequestRespondFromContext(c *gin.Context) (JoinRequestRespond, error) {
	var resp JoinRequestRespond
	if err := c.BindJSON(&resp); err != nil {
		return JoinRequestRespond{}, err
	}
	return resp, nil
}

func NoRouteHandler(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		SetHeaders(c)
		c.IndentedJSON(http_error.GetStatusCode(nil), nil)
	} else {
		err := http_error.NotFoundError
		c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
	}
}