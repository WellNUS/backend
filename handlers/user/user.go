package user

import (
	"wellnus/backend/references"
	"wellnus/backend/handlers/httpError"
	
	"fmt"
	"strconv"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"database/sql"
)

type User = references.User

func ReadUsers(rows *sql.Rows) ([]User, error) {
	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Gender, &user.Email, &user.UserRole, &user.PasswordHash); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func MergeUser(userMain User, userAdd User) (User, error) {
	userMain.ID = userAdd.ID
	if userMain.FirstName == "" {
		userMain.FirstName = userAdd.FirstName
	}
	if userMain.LastName == "" {
		userMain.LastName = userAdd.LastName
	}
	if userMain.Gender == "" {
		userMain.Gender = userAdd.Gender
	}
	if userMain.Email == "" {
		userMain.Email = userAdd.Email
	}
	if userMain.UserRole == "" {
		userMain.UserRole = userAdd.UserRole
	}
	if userMain.Password == "" {
		userMain.PasswordHash = userAdd.PasswordHash
	} else {
		var err error
		userMain.PasswordHash, err = argon2id.CreateHash(userMain.Password, argon2id.DefaultParams)
		userMain.Password = ""
		if err != nil { return User{}, err }	
	}
	return userMain, nil
}

func GetUser(db *sql.DB, id int64) (User, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM wn_user WHERE id = %d;", id))
	if err != nil { return User{}, err }
	defer rows.Close()

	users, err := ReadUsers(rows)
	if err != nil { return User{}, err}
	if len(users) == 0 { return User{}, httpError.NotFoundError }
	return users[0], nil
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM wn_user;")
	if err != nil { return nil, err }
	defer rows.Close()
	
	users, err := ReadUsers(rows)
	if err != nil { return nil, err}
	return users, nil
}

func AddUser(db *sql.DB, newUser User) (User, error) {
	var err error
	newUser.PasswordHash, err = argon2id.CreateHash(newUser.Password, argon2id.DefaultParams)
	newUser.Password = ""
	if err != nil { return User{}, err }

	_, err = db.Query(fmt.Sprintf(
		"INSERT INTO wn_user (first_name, last_name, gender, email, user_role, password_hash) VALUES ('%s', '%s', '%s', '%s', '%s', '%s');",
		newUser.FirstName,
		newUser.LastName,
		newUser.Gender,
		newUser.Email,
		newUser.UserRole,
		newUser.PasswordHash))
	if err != nil { return User{}, err }

	row, err := db.Query("SELECT last_value FROM wn_user_id_seq;")
	if err != nil { return User{}, err }
	defer row.Close()

	row.Next()
	if err := row.Scan(&newUser.ID); err != nil { return User{}, err }
	return newUser, nil
}

func DeleteUser(db *sql.DB, id int64) (User, error) {
	if _, err := db.Query(fmt.Sprintf("DELETE FROM wn_user WHERE id = %d", id)); err != nil {
		return User{}, err
	}
	return User{ID: id}, nil
}

func UpdateUser(db *sql.DB, updatedUser User, id int64) (User, error) {
	targetUser, err := GetUser(db, id)
	if err != nil { return User{}, err }

	updatedUser, err = MergeUser(updatedUser, targetUser)
	if err != nil { return User{}, err }

	query := fmt.Sprintf(
		"UPDATE wn_user SET first_name = '%s', last_name = '%s', gender = '%s', email = '%s', user_role = '%s', password_hash = '%s' WHERE id = %d;",
		updatedUser.FirstName,
		updatedUser.LastName,
		updatedUser.Gender,
		updatedUser.Email,
		updatedUser.UserRole,
		updatedUser.PasswordHash,
		id)
	if _, err := db.Query(query); err != nil {
		return User{}, err;
	}
	return updatedUser, nil;
}

// Handlers

func GetAllUsersHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		users, err := GetAllUsers(db)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), users)
	}
}

func GetUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		id, err := strconv.ParseInt(c.Param("id"), 0, 64)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		user, err := GetUser(db, id)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), user)
	}
}

func AddUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		var newUser User
		if err := c.BindJSON(&newUser); err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		newUser, err := AddUser(db, newUser)
		if err != nil { 
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		sid := fmt.Sprintf("%d", newUser.ID)
		c.SetCookie("id", sid, 1209600, "/", references.DOMAIN, false, true)
		c.IndentedJSON(httpError.GetStatusCode(err), newUser)
	}
}

func DeleteUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		id, err := strconv.ParseInt(c.Param("id"), 0, 64)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		strUserID, _ := c.Cookie("id")
		userID, _ := strconv.ParseInt(strUserID, 0, 64)
		if userID != id {
			err = httpError.UnauthorizedError
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		deletedUser, err := DeleteUser(db, id)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), deletedUser)
	}
}

func UpdateUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		id, err := strconv.ParseInt(c.Param("id"), 0, 64)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		strUserID, _ := c.Cookie("id")
		userID, _:= strconv.ParseInt(strUserID, 0, 64)
		if userID != id {
			err = httpError.UnauthorizedError
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		var updatedUser User
		if err := c.BindJSON(&updatedUser); err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		updatedUser, err = UpdateUser(db, updatedUser, id)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), updatedUser)
	}
}
