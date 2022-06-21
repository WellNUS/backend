package match

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/misc"
	"wellnus/backend/router/misc/http_error"

	"database/sql"
	"github.com/gin-gonic/gin"

)

func GetLoadedMatchRequestOfUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, err := misc.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		loadedMatchRequest, err := model.GetLoadedMatchRequestOfUser(db, userID)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), loadedMatchRequest)
	}
}

func AddMatchRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, err := misc.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		matchRequest, err := model.AddMatchRequest(db, userID)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), matchRequest)
	}
}

func DeleteMatchRequestOfUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, err := misc.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		matchRequest, err := model.DeleteMatchRequestOfUser(db, userID)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), matchRequest)
	}
}

func ForcePerformMatchingHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		model.PerformMatching(db)
		c.String(200, "Look at console")
	}
}