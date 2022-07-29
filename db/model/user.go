package model

import (
	"database/sql"
	"github.com/alexedwards/argon2id"
)

type User struct {
	ID 				int64 	`json:"id"`
	FirstName 		string 	`json:"first_name"`
	LastName 		string	`json:"last_name"`
	Gender			string 	`json:"gender"`
	Faculty			string 	`json:"faculty"`
	Email			string	`json:"email"`
	UserRole		string 	`json:"user_role"`
	Password		string 	`json:"password"`
	PasswordHash 	string	`json:"password_hash"`
}

type UserWithSessionKey struct {
	User 		User 	`json:"user"`
	SessionKey	string	`json:"session_key"`
}

type UserWithGroups struct {
	User 	User 	`json:"user"`
	Groups 	[]Group `json:"groups"`
}

func (user User) HashPassword() (User, error) {
	var err error
	user.PasswordHash, err = argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	user.Password = ""
	if err != nil { return User{}, err }
	return user, nil
}

func (user User) LoadLastUserID(db *sql.DB) (User, error) {
	row, err := db.Query("SELECT last_value FROM wn_user_id_seq;")
	if err != nil { return User{}, err }
	defer row.Close()

	row.Next()
	if err := row.Scan(&user.ID); err != nil { return User{}, err }
	return user, nil
}

func (userMain User) MergeUser(userAdd User) (User, error) {
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
	if userMain.Faculty == "" {
		userMain.Faculty = userAdd.Faculty
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

func (user1 User) Equal(user2 User) bool {
	return user1.ID == user2.ID &&
		user1.FirstName == user2.FirstName &&
		user1.LastName == user2.LastName &&
		user1.Gender == user2.Gender &&
		user1.Faculty == user2.Faculty &&
		user1.Email == user2.Email &&
		user1.UserRole == user2.UserRole &&
		user1.PasswordHash == user2.PasswordHash
}