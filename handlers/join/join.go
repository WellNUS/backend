package join

import (
	"wellnus/backend/config"
	"wellnus/backend/model"
	"wellnus/backend/db/query"
	"wellnus/backend/handlers/misc"

	"database/sql"
	"github.com/gin-gonic/gin"
)

type User = model.User
type Group = model.Group
type JoinRequest = model.JoinRequest
type LoadedJoinRequest = model.LoadedJoinRequest
type JoinRequestRespond = misc.JoinRequestRespond

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
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		userIDCookie, _ := misc.GetIDCookie(c)
		request := getRequestQuery(c)
		if request == REQUEST_RECEIVED {
			joinRequests, err := query.GetAllJoinRequestsReceivedOfUser(db, userIDCookie)
			if err != nil {
				c.IndentedJSON(misc.GetStatusCode(err), err.Error())
				return
			}
			c.IndentedJSON(misc.GetStatusCode(err), joinRequests)
		} else if request == REQUEST_SENT {
			joinRequests, err := query.GetAllJoinRequestsSentOfUser(db, userIDCookie)
			if err != nil {
				c.IndentedJSON(misc.GetStatusCode(err), err.Error())
				return
			}
			c.IndentedJSON(misc.GetStatusCode(err), joinRequests)
		} else {
			joinRequests, err := query.GetAllJoinRequestsOfUser(db, userIDCookie)
			if err != nil {
				c.IndentedJSON(misc.GetStatusCode(err), err.Error())
				return
			}
			c.IndentedJSON(misc.GetStatusCode(err), joinRequests)
		}
	}
}

func GetLoadedJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		joinRequestIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		loadedJoinRequest, err := query.GetLoadedJoinRequest(db, joinRequestIDParam)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), loadedJoinRequest)
	}
}

func AddJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		userIDCookie, err := misc.GetIDCookie(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err := misc.GetJoinRequestFromContext(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err = query.AddJoinRequest(db, joinRequest.GroupID, userIDCookie)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), joinRequest)
	}
}

func RespondJoinRequestHandler(db *sql.DB) func(*gin.Context) {	
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		userIDCookie, _ := misc.GetIDCookie(c)
		joinRequestIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		joinRequestRespond, err := misc.GetJoinRequestRespondFromContext(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		joinRequestRespond, err = query.RespondJoinRequest(db, joinRequestIDParam, userIDCookie, joinRequestRespond.Approve)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), joinRequestRespond)
	}	
}

func DeleteJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		userIDCookie, err := misc.GetIDCookie(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		joinRequestIDParam, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err := query.DeleteJoinRequest(db, joinRequestIDParam, userIDCookie)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), joinRequest)
	}
}