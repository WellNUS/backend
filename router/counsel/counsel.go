package counsel

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"

	"github.com/gin-gonic/gin"
	"database/sql"
)

func GetAllCounselRequestsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		topics, _ := c.GetQueryArray("topic")
		counselRequests, err := model.GetAllCounselRequests(db, topics, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), counselRequests)
	}
}

func GetCounselRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userIDCookie, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		counselRequest, err := model.GetCounselRequest(db, userIDParam, userIDCookie)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), counselRequest)
	}
}

func AddUpdateCounselRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		counselRequest, err := http_helper.GetCounselRequestFromContext(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		counselRequest, err = model.AddUpdateCounselRequest(db, counselRequest, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), counselRequest)
	}
}

func DeleteCounselRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {	
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		counselRequest, err := model.DeleteCounselRequest(db, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), counselRequest)
	}
}

func AcceptCounselRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)
		
		userIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		groupWithUsers, err := model.AcceptCounselRequest(db, userIDParam, userIDCookie)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), groupWithUsers)
	}
}
