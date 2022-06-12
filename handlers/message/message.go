package message

import (
	"wellnus/backend/db/query"
	"wellnus/backend/handlers/misc"

	"database/sql"

	"github.com/gin-gonic/gin"
)

func GetAllLoadedMessagesOfGroup(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		groupID, err := misc.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		loadedMessages, err := query.GetAllLoadedMessagesOfGroup(db, groupID)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(misc.GetStatusCode(err), loadedMessages)
	}
}