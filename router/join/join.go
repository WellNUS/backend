package join

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"

	"database/sql"
	"github.com/gin-gonic/gin"
)

const (
	REQUEST_RECEIVED = 0
	REQUEST_SENT 	= 1
	REQUEST_BOTH 	= 2
)

// Helper functions

func getRequestQuery(c *gin.Context) int {
	if s := c.Query("request"); s == "RECEIVED" {
		return REQUEST_RECEIVED
	} else if s == "SENT" {
		return REQUEST_SENT
	} else {
		return REQUEST_BOTH
	}
}

// Main functions

func GetAllLoadedJoinRequestsHandler(db *sql.DB) func(*gin.Context){
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		request := getRequestQuery(c)
		if request == REQUEST_RECEIVED {
			joinRequests, err := model.GetAllLoadedJoinRequestsReceivedOfUser(db, userID)
			if err != nil {
				c.JSON(http_error.GetStatusCode(err), err.Error())
				return
			}
			c.JSON(http_error.GetStatusCode(err), joinRequests)
		} else if request == REQUEST_SENT {
			joinRequests, err := model.GetAllLoadedJoinRequestsSentOfUser(db, userID)
			if err != nil {
				c.JSON(http_error.GetStatusCode(err), err.Error())
				return
			}
			c.JSON(http_error.GetStatusCode(err), joinRequests)
		} else {
			joinRequests, err := model.GetAllLoadedJoinRequestsOfUser(db, userID)
			if err != nil {
				c.JSON(http_error.GetStatusCode(err), err.Error())
				return
			}
			c.JSON(http_error.GetStatusCode(err), joinRequests)
		}
	}
}

func GetLoadedJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		joinRequestIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		loadedJoinRequest, err := model.GetLoadedJoinRequest(db, joinRequestIDParam)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), loadedJoinRequest)
	}
}

func AddJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err := http_helper.GetJoinRequestFromContext(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err = model.AddJoinRequest(db, joinRequest.GroupID, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), joinRequest)
	}
}

func RespondJoinRequestHandler(db *sql.DB) func(*gin.Context) {	
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		joinRequestIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		joinRequestRespond, err := http_helper.GetJoinRequestRespondFromContext(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		joinRequestRespond, err = model.RespondJoinRequest(db, joinRequestIDParam, userID, joinRequestRespond.Approve)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), joinRequestRespond)
	}	
}

func DeleteJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		joinRequestIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err := model.DeleteJoinRequest(db, joinRequestIDParam, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), joinRequest)
	}
}