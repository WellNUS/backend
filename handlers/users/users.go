package users

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

func getUser(db *sql.DB, id int64) (User, error) {
	row, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %d;", id))
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

func getAllUsers(db *sql.DB) ([]User, error) {
	users := make([]User, 0)
	rows, err := db.Query("SELECT * FROM users;")
	if err != nil { return nil, err }
	defer rows.Close()
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Gender, &user.Email, &user.PasswordHash); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func addUser(db *sql.DB, newUser User) (User, error) {
	var err error
	newUser.PasswordHash, err = argon2id.CreateHash(newUser.Password, argon2id.DefaultParams)
	newUser.Password = ""
	if err != nil { return User{}, err }
	db.Query(fmt.Sprintf(
		"INSERT INTO users (first_name, last_name, gender, email, password_hash) VALUES ('%s', '%s', '%s', '%s', '%s');",
		newUser.FirstName,
		newUser.LastName,
		newUser.Gender,
		newUser.Email,
		newUser.PasswordHash))
	row, err := db.Query("SELECT last_value FROM users_id_seq;")
	if err != nil { return User{}, err }
	row.Next()
	if err := row.Scan(&newUser.ID); err != nil { return User{}, err }
	return newUser, nil
}

func deleteUser(db *sql.DB, id int64) (User, error) {
	if _, err := db.Query(fmt.Sprintf("DELETE FROM users WHERE id = %d", id)); err != nil {
		return User{}, err
	}
	return User{ID: id}, nil
}

func updateUser(db *sql.DB, updatedUser User, id int64) (User, error) {
	targetUser, err := getUser(db, id)
	if err != nil {
		return User{}, err
	}
	updatedUser, err = mergeUser(updatedUser, targetUser)
	if err != nil {
		return User{}, err
	}
	query := fmt.Sprintf(
		"UPDATE users SET first_name = '%s', last_name = '%s', gender = '%s', email = '%s', password_hash = '%s' WHERE id = %d;",
		updatedUser.FirstName,
		updatedUser.LastName,
		updatedUser.Gender,
		updatedUser.Email,
		updatedUser.PasswordHash,
		id)
	if _, err := db.Query(query); err != nil {
		return User{}, err;
	}
	return updatedUser, nil;
}

func mergeUser(userMain User, userAdd User) (User, error) {
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

// Handlers

func GetAllUsersHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", references.FRONTEND_URL)
    	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		users, err := getAllUsers(db)
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
		user, err := getUser(db, id)
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
		newUser, err := addUser(db, newUser)
		if err != nil { 
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
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
		deletedUser, err := deleteUser(db, id)
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
		updatedUser, err = updateUser(db, updatedUser, id)
		if err != nil {
			c.IndentedJSON(httpError.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(httpError.GetStatusCode(err), updatedUser)
	}
}
