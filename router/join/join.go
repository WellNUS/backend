package join

import (
	"wellnus/backend/config"
	"wellnus/backend/db/model"
	"wellnus/backend/router/misc"
	"wellnus/backend/router/misc/http_error"

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

func GetAllJoinRequestsHandler(db *sql.DB) func(*gin.Context){
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")
		userIDCookie, _ := misc.GetIDCookie(c)
		request := getRequestQuery(c)
		if request == REQUEST_RECEIVED {
			joinRequests, err := model.GetAllJoinRequestsReceivedOfUser(db, userIDCookie)
			if err != nil {
				c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
				return
			}
			c.IndentedJSON(http_error.GetStatusCode(err), joinRequests)
		} else if request == REQUEST_SENT {
			joinRequests, err := model.GetAllJoinRequestsSentOfUser(db, userIDCookie)
			if err != nil {
				c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
				return
			}
			c.IndentedJSON(http_error.GetStatusCode(err), joinRequests)
		} else {
			joinRequests, err := model.GetAllJoinRequestsOfUser(db, userIDCookie)
			if err != nil {
				c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
				return
			}
			c.IndentedJSON(http_error.GetStatusCode(err), joinRequests)
		}
	}
}

func GetLoadedJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")
		joinRequestIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		loadedJoinRequest, err := model.GetLoadedJoinRequest(db, joinRequestIDParam)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), loadedJoinRequest)
	}
}

func AddJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")
		userIDCookie, err := misc.GetIDCookie(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err := misc.GetJoinRequestFromContext(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err = model.AddJoinRequest(db, joinRequest.GroupID, userIDCookie)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), joinRequest)
	}
}

func RespondJoinRequestHandler(db *sql.DB) func(*gin.Context) {	
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")
		userIDCookie, _ := misc.GetIDCookie(c)
		joinRequestIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		joinRequestRespond, err := misc.GetJoinRequestRespondFromContext(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		joinRequestRespond, err = model.RespondJoinRequest(db, joinRequestIDParam, userIDCookie, joinRequestRespond.Approve)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), joinRequestRespond)
	}	
}

func DeleteJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")
		userIDCookie, err := misc.GetIDCookie(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		joinRequestIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err := model.DeleteJoinRequest(db, joinRequestIDParam, userIDCookie)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), joinRequest)
	}
}