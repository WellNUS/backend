package ws

import (
	"github.com/gin-gonic/gin"
)

func ConnectToWSHandler(wsHub *Hub) func(*gin.Context) {
	return func(c *gin.Context) {
		ServeWs(wsHub, c.Writer, c.Request)
	}
}