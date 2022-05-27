package session

import (
	"wellnus/backend/references"
	"wellnus/backend/handlers/httpError"

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
	row, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE email = '%s';", email))
	if err != nil { return User{}, err }
	if row.Next() {
		var user User
		if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Gender, &user.Email, &user.PasswordHash); err != nil {
			return User{}, err
		}
		return user, nil
	}
	return User{}, httpError.NotFoundError
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
			id := fmt.Sprintf("%d",storedUser.ID)
			c.SetCookie("id", id, 1209600, "/", references.DOMAIN, false, true)
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
