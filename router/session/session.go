package session

import (
	"wellnus/backend/router/misc"
	"wellnus/backend/router/misc/http_error"
	"wellnus/backend/db/model"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"database/sql"
)

type User = model.User
type Resp = model.Resp

// Main function
func LoginHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		misc.SetHeaders(c)

		loginUser, err := misc.GetUserFromContext(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		storedUser, err := model.FindUser(db, loginUser.Email)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		match, err := argon2id.ComparePasswordAndHash(loginUser.Password, storedUser.PasswordHash)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		if match {
			misc.SetIDCookie(c, storedUser.ID)
			c.IndentedJSON(http_error.GetStatusCode(err), Resp{ LoggedIn: true, User: storedUser })
		} else {
			misc.RemoveIDCookie(c)
			c.IndentedJSON(http_error.GetStatusCode(err), Resp{ LoggedIn: false, User: User{}})
		}
	} 
}

func LogoutHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		misc.SetHeaders(c)

		misc.RemoveIDCookie(c)
		c.IndentedJSON(http_error.GetStatusCode(nil), Resp{ LoggedIn: false, User: User{}})
	}
}