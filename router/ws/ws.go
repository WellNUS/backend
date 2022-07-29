package ws

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"
	"fmt"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func ConnectToWSHandler(wsHub *Hub, db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		groupID, err := http_helper.GetIDParams(c)
		if err != nil {
			fmt.Printf("An error occured when retrieving group ID params. %v \n", err)
			return
		}
		userID, err := http_helper.GetUserIDFromSessionHeader(db, c)
		if err != nil {
			fmt.Printf("An error occured when retrieving user ID cookies. %v \n", err)
			return
		}
		isMember, err := model.IsUserInGroup(db, userID, groupID)
		if err != nil {
			fmt.Printf("An error occured when checking if user is in group. %v \n", err)
			return
		}
		if !isMember {
			err = http_error.UnauthorizedError
			fmt.Printf("User is not part of group. %v \n", err)
			return
		}
		ServeWs(wsHub, c.Writer, c.Request, userID, groupID)
	}
}