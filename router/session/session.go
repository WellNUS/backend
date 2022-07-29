package session

import (
	"wellnus/backend/config"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"
	"wellnus/backend/db/model"

	"net/http"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"database/sql"
)

// JWT vs Cookie authentication
// App currently uses cookie authentication. Sessions are stored in database
// JWT uses one secret key across the whole app. When request sent with JWT, secret key is decripted to get userID. No session required in database

type User = model.User
type SessionResponse = model.SessionResponse

// Helper function
func CreateNewSessionCookie(db *sql.DB, c *gin.Context, userID int64) error {
	newSessionKey, err := model.CreateNewSession(db, userID)
	if err != nil { return err }
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("session_key", newSessionKey, 1209600, "/", config.COOKIE_ADDRESS, true, true)
	return nil
}

func RemoveSessionCookie(db *sql.DB, c *gin.Context) error {
	sessionKey, _ := c.Cookie("session_key")
	if err := model.DeleteSessionWithSessionKey(db, sessionKey); err != nil { return err }
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("session_key", "", -1, "/", config.COOKIE_ADDRESS, true, true)
	return nil
}

// Main function
func LoginHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		loginUser, err := http_helper.GetUserFromContext(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		storedUser, err := model.FindUser(db, loginUser.Email)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		match, err := argon2id.ComparePasswordAndHash(loginUser.Password, storedUser.PasswordHash)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		if match {
			// err = CreateNewSessionCookie(db, c, storedUser.ID)
			sessionKey, err := model.CreateNewSession(db, storedUser.ID)
			if err != nil {
				c.JSON(http_error.GetStatusCode(err), err.Error())
				return
			}
			c.JSON(http_error.GetStatusCode(err), SessionResponse{ LoggedIn: true, User: storedUser, SessionKey: sessionKey })
		} else {
			RemoveSessionCookie(db, c)
			c.JSON(http_error.GetStatusCode(err), SessionResponse{ LoggedIn: false })
		}
	} 
}

func LogoutHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		err := model.DeleteSessionWithSessionKey(db, c.Request.Header.Get("session_key"))
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(nil), SessionResponse{ LoggedIn: false, User: User{}})
	}
}