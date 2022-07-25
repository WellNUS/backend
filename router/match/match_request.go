package match

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"

	"database/sql"
	"github.com/gin-gonic/gin"

)

func GetMatchRequestCount(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		count, err := model.GetMatchRequestCount(db)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), count)
	}
}

func GetLoadedMatchRequestOfUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		paramID, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
		}
		if userID != paramID {
			err = http_error.UnauthorizedError
			c.JSON(http_error.GetStatusCode(err), err.Error())
		}

		loadedMatchRequest, err := model.GetLoadedMatchRequestOfUser(db, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), loadedMatchRequest)
	}
}

func AddMatchRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		matchRequest, err := model.AddMatchRequest(db, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), matchRequest)
	}
}

func DeleteMatchRequestOfUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)
		
		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		matchRequest, err := model.DeleteMatchRequestOfUser(db, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), matchRequest)
	}
}