package ws

import (
	"wellnus/backend/handlers/misc"
	"wellnus/backend/db/query"

	"fmt"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func ConnectToWSHandler(wsHub *Hub, db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		groupID, err := misc.GetIDParams(c)
		if err != nil {
			fmt.Printf("An error occured when retrieving group ID params. %v \n", err)
			return
		}
		userID, err := misc.GetIDCookie(c)
		if err != nil {
			fmt.Printf("An error occured when retrieving user ID cookies. %v \n", err)
			return
		}
		isMember, err := query.IsUserInGroup(db, userID, groupID)
		if err != nil {
			fmt.Printf("An error occured when checking if user is in group. %v \n", err)
			return
		}
		if !isMember {
			err = misc.UnauthorizedError
			fmt.Printf("User is not part of group. %v \n", err)
			return
		}
		ServeWs(wsHub, c.Writer, c.Request, userID, groupID)
	}
}