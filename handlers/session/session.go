package session

import (
	"wellnus/backend/references"
	"wellnus/backend/handlers/httpError"
	"wellnus/backend/handlers/user"

	"fmt"
	// "strconv"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"database/sql"
)

type User = references.User

type Resp struct {
	LoggedIn 	bool `json:"logged_in"`
	User	 	User `json:"user"`
}

func findUser(db *sql.DB, email string) (User, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM wn_user WHERE email = '%s';", email))
	if err != nil { return User{}, err }
	users, err := user.ReadUsers(rows)
	if err != nil { return User{}, err}
	if len(users) == 0 { return User{}, httpError.NotFoundError }
	return users[0], nil
}

func LoginHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		var loginUser User
		if err := c.BindJSON(&loginUser); err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		storedUser, err := findUser(db, loginUser.Email)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		match, err := argon2id.ComparePasswordAndHash(loginUser.Password, storedUser.PasswordHash)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		if match {
			sid := fmt.Sprintf("%d",storedUser.ID)
			c.SetCookie("id", sid, 1209600, "/", references.DOMAIN, false, true)
			c.IndentedJSON(httpError.GetStatusCode(err), Resp{ LoggedIn: true, User: storedUser })
		} else {
			c.SetCookie("id", "", -1, "/", references.DOMAIN, false, true)
			c.IndentedJSON(httpError.GetStatusCode(err),Resp{ LoggedIn: false, User: User{}})
		}
	} 
}

func LogoutHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var err error
		c.SetCookie("id", "", -1, "/", references.DOMAIN, false, true)
		c.IndentedJSON(httpError.GetStatusCode(err), Resp{ LoggedIn: false, User: User{}})
	}
}
