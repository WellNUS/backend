package message

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"

	"database/sql"

	"github.com/gin-gonic/gin"
)

type DirectMessagePayload = model.DirectMessagePayload

func GetDirectMessagesChunkHandler(db *sql.DB) func(*gin.Context) {
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

		limit, _ := http_helper.GetLimitQuery(c)
		latestTime, err := http_helper.GetLatestQuery(c)
		if err != nil { 
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		
		directMessagesChunk, err := model.GetDirectMessagesChunk(db, userIDCookie, userIDParam, latestTime, limit)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), directMessagesChunk)
	}
}