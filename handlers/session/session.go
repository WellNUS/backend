package session

import (
	"wellnus/backend/config"
	"wellnus/backend/model"
	"wellnus/backend/handlers/misc"
	"wellnus/backend/db/query"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"database/sql"
)

type User = model.User

type Resp struct {
	LoggedIn 	bool `json:"logged_in"`
	User	 	User `json:"user"`
}

// Main function
func LoginHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		loginUser, err := misc.GetUserFromContext(c)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}

		storedUser, err := query.FindUser(db, loginUser.Email)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}

		match, err := argon2id.ComparePasswordAndHash(loginUser.Password, storedUser.PasswordHash)
		if err != nil {
			c.IndentedJSON(misc.GetStatusCode(err), err.Error())
			return
		}
		if match {
			misc.SetIDCookie(c, storedUser.ID)
			c.IndentedJSON(misc.GetStatusCode(err), Resp{ LoggedIn: true, User: storedUser })
		} else {
			misc.RemoveIDCookie(c)
			c.IndentedJSON(misc.GetStatusCode(err), Resp{ LoggedIn: false, User: User{}})
		}
	} 
}

func LogoutHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		misc.RemoveIDCookie(c)
		c.IndentedJSON(misc.GetStatusCode(nil), Resp{ LoggedIn: false, User: User{}})
	}
}