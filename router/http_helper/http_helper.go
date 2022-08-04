package http_helper

import (
	"wellnus/backend/db/model"
	"wellnus/backend/config"
	"wellnus/backend/router/http_helper/http_error"
	"strconv"
	"time"
	"log"
	"database/sql"
	"github.com/gin-gonic/gin"
)

type User = model.User
type Group = model.Group
type JoinRequest = model.JoinRequest
type JoinRequestRespond = model.JoinRequestRespond
type MatchSetting = model.MatchSetting
type ProviderSetting = model.ProviderSetting
type CounselRequest = model.CounselRequest
type UserIDBody = model.UserIDBody
type Event = model.Event
type Booking = model.Booking
type BookingRespond = model.BookingRespond

func SetHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", config.FRONTEND_ADDRESS)
	c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
	c.Header("Access-Control-Allow-Credentials", "true")
}

func GetIDParams(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil { return 0, http_error.NotFoundError }
	return id, nil
}

func GetLatestQuery(c *gin.Context) (time.Time, error) {
	// RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	if stime := c.Query("latest"); stime == "" {
		return time.Now(), nil
	} else {
		return time.Parse(time.RFC3339Nano, stime)
	}
}

func GetLimitQuery(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Query("limit"), 0, 64)
}

func GetUserIDFromSessionCookie(db *sql.DB, c *gin.Context) (int64, error) {
	sessionKey, err := c.Cookie("session_key")
	if err != nil { 
		log.Printf("Error while getting UserID from cookie: %v", err)
		return 0, http_error.UnauthorizedError
	}
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

func GetProviderSettingFromContext(c *gin.Context) (ProviderSetting, error) {
	var providerSetting ProviderSetting
	if err := c.BindJSON(&providerSetting); err != nil {
		return ProviderSetting{}, err
	}
	return providerSetting, nil
}

func GetMatchSettingFromContext(c *gin.Context) (MatchSetting, error) {
	var matchSetting MatchSetting
	if err := c.BindJSON(&matchSetting); err != nil {
		return MatchSetting{}, err
	}
	return matchSetting, nil
}

func GetCounselRequestFromContext(c *gin.Context) (CounselRequest, error) {
	var counselRequest CounselRequest
	if err := c.BindJSON(&counselRequest); err != nil {
		return CounselRequest{}, err
	}
	return counselRequest, nil
}

func GetEventFromContext(c *gin.Context) (Event, error) {
	var event Event
	if err := c.BindJSON(&event); err != nil {
		return Event{}, err
	}
	return event, nil
}

func GetUserIDFromContext(c *gin.Context) (int64, error) {
	var userIDBody UserIDBody
	if err := c.BindJSON(&userIDBody); err != nil {
		return 0, err
	}
	return userIDBody.UserID, nil
}

func GetBookingFromContext(c *gin.Context) (Booking, error) {
	var booking Booking
	if err := c.BindJSON(&booking); err != nil {
		return Booking{}, err
	}
	return booking, nil
}

func GetBookingRespondFromContext(c *gin.Context) (BookingRespond, error) {
	var bookingRespond BookingRespond
	if err := c.BindJSON(&bookingRespond); err != nil {
		return BookingRespond{}, err
	}
	return bookingRespond, nil
}

func NoRouteHandler(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		SetHeaders(c)
		c.JSON(http_error.GetStatusCode(nil), nil)
	} else {
		err := http_error.NotFoundError
		c.JSON(http_error.GetStatusCode(err), err.Error())
	}
}