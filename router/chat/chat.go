package chat

import (
	"wellnus/backend/config"
	"wellnus/backend/db/model"
	"wellnus/backend/router/misc"
	"wellnus/backend/router/misc/http_error"

	"time"
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessagePayload = model.MessagePayload

func getLatestQuery(c *gin.Context) (time.Time, error) {
	// RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	if stime := c.Query("latest"); stime == "" {
		return time.Now(), nil
	} else {
		return time.Parse(time.RFC3339Nano, stime)
	}
}

func getLimitQuery(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Query("limit"), 0, 64)
}

func GetMessagesChunkOfGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
		c.Header("Access-Control-Allow-Methods", "PATCH, POST, GET, DELETE, OPTIONS")

		groupID, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		limit, _ := getLimitQuery(c)
		latestTime, err := getLatestQuery(c)
		if err != nil { 
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		
		messagesChunk, err := model.GetMessagesChunkOfGroupCustomise(db, groupID, latestTime, limit)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), messagesChunk)
	}
}