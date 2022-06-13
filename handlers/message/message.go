package message

import (
	"wellnus/backend/db/query"
	"wellnus/backend/db/model"
	"wellnus/backend/handlers/misc"

	"time"
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LoadedMessage = model.LoadedMessage
type LoadedMessagesPacket = model.LoadedMessagesPacket

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

func GetLoadedMessagesPacketOfGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		groupID, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}

		limit, _ := getLimitQuery(c)
		latestTime, err := getLatestQuery(c)
		if err != nil { 
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		
		loadedMessages, err := query.GetLoadedMessagesOfGroupCustomise(db, groupID, latestTime, limit)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}

		var loadedMessagesPacket LoadedMessagesPacket
		if l := len(loadedMessages); l > 0 {
			loadedMessagesPacket = LoadedMessagesPacket{
				EarliestTime: loadedMessages[0].Message.TimeAdded,
				LatestTime: loadedMessages[len(loadedMessages) - 1].Message.TimeAdded,
				LoadedMessages: loadedMessages,
			}
		} else {
			loadedMessagesPacket = LoadedMessagesPacket{LoadedMessages: loadedMessages}
		}
		c.IndentedJSON(misc.GetStatusCode(err), loadedMessagesPacket)
	}
}