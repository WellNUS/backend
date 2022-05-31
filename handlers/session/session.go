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

// Helper function
func findUser(db *sql.DB, email string) (User, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM wn_user WHERE email = '%s';", email))
	if err != nil { return User{}, err }
	users, err := readUsers(rows)
	if err != nil { return User{}, err}
	if len(users) == 0 { return User{}, httpError.NotFoundError }
	return users[0], nil
}

func readUsers(rows *sql.Rows) ([]User, error) {
	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Gender, &user.Faculty, &user.Email, &user.UserRole, &user.PasswordHash); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func getUserFromContext(c *gin.Context) (User, error) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		return User{}, err
	}
	return user, nil
}

func setIDCookie(c *gin.Context, id int64) {
	sid := fmt.Sprintf("%d",id)
	c.SetCookie("id", sid, 1209600, "/", references.DOMAIN, false, true)
}

func removeIDCookie(c *gin.Context) {
	c.SetCookie("id", "", -1, "/", references.DOMAIN, false, true)
}

// Main function
func LoginHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		loginUser, err := getUserFromContext(c)
		if err != nil {
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
			setIDCookie(c, storedUser.ID)
			c.IndentedJSON(httpError.GetStatusCode(err), Resp{ LoggedIn: true, User: storedUser })
		} else {
			removeIDCookie(c)
			c.IndentedJSON(httpError.GetStatusCode(err), Resp{ LoggedIn: false, User: User{}})
		}
	} 
}

func LogoutHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		removeIDCookie(c)
		c.IndentedJSON(httpError.GetStatusCode(nil), Resp{ LoggedIn: false, User: User{}})
	}
}