package join

import (
	"wellnus/backend/references"
	"wellnus/backend/handlers/httpError"
	"database/sql"
	"github.com/gin-gonic/gin"

	"strconv"
)

type User = references.User
type Group = references.Group
type JoinRequest = references.JoinRequest
type JoinRequestWithGroup = references.JoinRequestWithGroup
type JoinRequestRespond struct {
	Approve bool `json:"approve"`
}

const (
	REQUEST_RECEIVED = 0
	REQUEST_SENT 	= 1
	REQUEST_BOTH 	= 2
)

// Helper functions

func getIDParams(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil { return 0, httpError.NotFoundError }
	return id, nil
}

func getJoinRequestFromContext(c *gin.Context) (JoinRequest, error) {
	var joinRequest JoinRequest
	if err := c.BindJSON(&joinRequest); err != nil {
		return JoinRequest{}, err
	}
	return joinRequest, nil
}

func getIDCookie(c *gin.Context) (int64, error) {
	strUserID, err := c.Cookie("id")
	if err != nil { return 0, httpError.UnauthorizedError }
	userID, err := strconv.ParseInt(strUserID, 0, 64)
	if err != nil { return 0, err }
	return userID, nil
}

func getRequestQuery(c *gin.Context) int {
	if s := c.Query("request"); s == "RECEIVED" {
		return REQUEST_RECEIVED
	} else if s == "SENT" {
		return REQUEST_SENT
	} else {
		return REQUEST_BOTH
	}
}

func getJoinRequestRespondFromContext(c *gin.Context) (JoinRequestRespond, error) {
	var resp JoinRequestRespond
	if err := c.BindJSON(&resp); err != nil {
		return JoinRequestRespond{}, err
	}
	return resp, nil
}	

// Main functions

func GetAllJoinRequestsHandler(db *sql.DB) func(*gin.Context){
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		userIDCookie, _ := getIDCookie(c)
		request := getRequestQuery(c)
		if request == REQUEST_RECEIVED {
			joinRequests, err := GetAllJoinRequestsReceived(db, userIDCookie)
			if err != nil {
				c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
				return
			}
			c.IndentedJSON(httpError.GetStatusCode(err), joinRequests)
		} else if request == REQUEST_SENT {
			joinRequests, err := GetAllJoinRequestsSent(db, userIDCookie)
			if err != nil {
				c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
				return
			}
			c.IndentedJSON(httpError.GetStatusCode(err), joinRequests)
		} else {
			joinRequests, err := GetAllJoinRequests(db, userIDCookie)
			if err != nil {
				c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
				return
			}
			c.IndentedJSON(httpError.GetStatusCode(err), joinRequests)
		}
	}
}

func GetJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		joinRequestIDParam, err := getIDParams(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err := GetJoinRequest(db, joinRequestIDParam)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), joinRequest)
	}
}

func AddJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		userIDCookie, err := getIDCookie(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err := getJoinRequestFromContext(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err = AddJoinRequest(db, joinRequest.GroupID, userIDCookie)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), joinRequest)
	}
}

func RespondJoinRequestHandler(db *sql.DB) func(*gin.Context) {	
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		userIDCookie, _ := getIDCookie(c)
		joinRequestIDParam, err := getIDParams(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		resp, err := getJoinRequestRespondFromContext(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err := RespondJoinRequest(db, joinRequestIDParam, userIDCookie, resp.Approve)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), joinRequest)
	}	
}

func DeleteJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		userIDCookie, err := getIDCookie(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		joinRequestIDParam, err := getIDParams(c)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		joinRequest, err := DeleteJoinRequest(db, joinRequestIDParam, userIDCookie)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), joinRequest)
	}
}