package session

import (
	"wellnus/backend/config"
	"wellnus/backend/router/misc"
	"wellnus/backend/router/misc/http_error"
	"wellnus/backend/db/model"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"database/sql"
)

type User = model.User
type SessionResponse = model.SessionResponse

// Helper function
func CreateNewSessionCookie(db *sql.DB, c *gin.Context, userID int64) error {
	newSessionKey, err := model.CreateNewSession(db, userID)
	if err != nil { return err }
	c.SetCookie("session_key", newSessionKey, 1209600, "/", config.DOMAIN, false, true)
	return nil
}

func RemoveSessionCookie(db *sql.DB, c *gin.Context) error {
	sessionKey, _ := c.Cookie("session_key")
	if err := model.DeleteSessionWithSessionKey(db, sessionKey); err != nil { return err }
	c.SetCookie("session_key", "", -1, "/", config.DOMAIN, false, true)
	return nil
}

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
			err = CreateNewSessionCookie(db, c, storedUser.ID)
			if err != nil {
				c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
				return
			}
			c.IndentedJSON(http_error.GetStatusCode(err), SessionResponse{ LoggedIn: true, User: storedUser })
		} else {
			RemoveSessionCookie(db, c)
			c.IndentedJSON(http_error.GetStatusCode(err), SessionResponse{ LoggedIn: false, User: User{}})
		}
	} 
}

func LogoutHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		misc.SetHeaders(c)

		err := RemoveSessionCookie(db, c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(nil), SessionResponse{ LoggedIn: false, User: User{}})
	}
}