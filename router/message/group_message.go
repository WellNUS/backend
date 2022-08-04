package message

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"

	"database/sql"

	"github.com/gin-gonic/gin"
)

type GroupMessagePayload = model.GroupMessagePayload

func GetGroupMessagesChunkHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		groupID, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		inGroup, err := model.IsUserInGroup(db, userID, groupID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		if !inGroup {
			err = http_error.UnauthorizedError
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		limit, _ := http_helper.GetLimitQuery(c)
		latestTime, err := http_helper.GetLatestQuery(c)
		if err != nil { 
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		
		groupMessagesChunk, err := model.GetGroupMessagesChunk(db, groupID, latestTime, limit)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), groupMessagesChunk)
	}
}