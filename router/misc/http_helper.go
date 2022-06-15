package misc

import (
	"wellnus/backend/db/model"
	"wellnus/backend/config"
	"wellnus/backend/router/misc/http_error"
	"strconv"
	"fmt"
	"github.com/gin-gonic/gin"
)

type User = model.User
type Group = model.Group
type JoinRequest = model.JoinRequest
type JoinRequestRespond = model.JoinRequestRespond

func GetIDParams(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil { return 0, http_error.NotFoundError }
	return id, nil
}

func GetIDCookie(c *gin.Context) (int64, error) {
	strUserID, err := c.Cookie("id")
	if err != nil { return 0, http_error.UnauthorizedError }
	userID, err := strconv.ParseInt(strUserID, 0, 64)
	if err != nil { return 0, err }
	return userID, nil
}

func SetIDCookie(c *gin.Context, userID int64) {
	c.SetCookie("id", fmt.Sprintf("%d", userID), 1209600, "/", config.DOMAIN, false, true)
}

func RemoveIDCookie(c *gin.Context) {
	c.SetCookie("id", "", -1, "/", config.DOMAIN, false, true)
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