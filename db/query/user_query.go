package query

import (
	"wellnus/backend/handlers/misc"
	"wellnus/backend/db/model"
	
	"fmt"
	"github.com/alexedwards/argon2id"
	"database/sql"
)

type User = model.User
type UserWithGroups = model.UserWithGroups

//Helper functions

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

func mergeUser(userMain User, userAdd User) (User, error) {
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

func hashPassword(user User) (User, error) {
	var err error
	user.PasswordHash, err = argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	user.Password = ""
	if err != nil { return User{}, err }
	return user, nil
}

func getUser(db *sql.DB, id int64) (User, error) {
	query := fmt.Sprintf("SELECT * FROM wn_user WHERE id = %d;", id)
	rows, err := db.Query(query)
	if err != nil { return User{}, err }
	defer rows.Close()

	users, err := readUsers(rows)
	if err != nil { return User{}, err}
	if len(users) == 0 { return User{}, misc.NotFoundError }
	return users[0], nil
}

func loadLastUserID(db *sql.DB, user User) (User, error) {
	row, err := db.Query("SELECT last_value FROM wn_user_id_seq;")
	if err != nil { return User{}, err }
	defer row.Close()

	row.Next()
	if err := row.Scan(&user.ID); err != nil { return User{}, err }
	return user, nil
}

// Main functions

func GetUserWithGroups(db *sql.DB, userID int64) (UserWithGroups, error) {
	user, err := getUser(db, userID)
	if err != nil { return UserWithGroups{}, err }
	groups, err := GetAllGroupsOfUser(db, userID)
	if err != nil { return UserWithGroups{}, err }
	return UserWithGroups{ User: user, Groups: groups}, nil
}

func GetAllUsersOfGroup(db *sql.DB, groupID int64) ([]User, error) {
	query := fmt.Sprintf(
		`SELECT 
			wn_user.id,
			wn_user.first_name,
			wn_user.last_name,
			wn_user.gender,
			wn_user.faculty,
			wn_user.email,
			wn_user.user_role,
			wn_user.password_hash
		FROM wn_user_group JOIN wn_user 
		ON wn_user_group.user_id = wn_user.id 
		WHERE wn_user_group.group_id = %d`, 
		groupID)
	rows, err := db.Query(query)
	if err != nil { return nil, err }
	users, err := readUsers(rows)
	if err != nil { return nil, err }
	return users, nil
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM wn_user;")
	if err != nil { return nil, err }
	defer rows.Close()
	
	users, err := readUsers(rows)
	if err != nil { return nil, err}
	return users, nil
}

func AddUser(db *sql.DB, newUser User) (User, error) {
	newUser, err := hashPassword(newUser)
	if err != nil { return User{}, err }
	query := fmt.Sprintf(
		`INSERT INTO wn_user (
			first_name, 
			last_name, 
			gender, 
			faculty, 
			email, 
			user_role, 
			password_hash
		) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');`,
		newUser.FirstName,
		newUser.LastName,
		newUser.Gender,
		newUser.Faculty,
		newUser.Email,
		newUser.UserRole,
		newUser.PasswordHash)
	_, err = db.Query(query)
	if err != nil { return User{}, err }

	newUser, err = loadLastUserID(db, newUser)
	if err != nil { return User{}, err }
	return newUser, nil
}

func DeleteUser(db *sql.DB, id int64) (User, error) {
	if _, err := db.Query(fmt.Sprintf("DELETE FROM wn_user WHERE id = %d", id)); err != nil {
		return User{}, err
	}
	return User{ID: id}, nil
}

func UpdateUser(db *sql.DB, updatedUser User, id int64) (User, error) {
	targetUser, err := getUser(db, id)
	if err != nil { return User{}, err }

	updatedUser, err = mergeUser(updatedUser, targetUser)
	if err != nil { return User{}, err }

	query := fmt.Sprintf(
		`UPDATE wn_user SET 
			first_name = '%s', 
			last_name = '%s', 
			gender = '%s', 
			faculty='%s', 
			email = '%s', 
			user_role = '%s', 
			password_hash = '%s' 
		WHERE id = %d;`,
		updatedUser.FirstName,
		updatedUser.LastName,
		updatedUser.Gender,
		updatedUser.Faculty,
		updatedUser.Email,
		updatedUser.UserRole,
		updatedUser.PasswordHash,
		id)
	if _, err := db.Query(query); err != nil {
		return User{}, err;
	}
	return updatedUser, nil;
}

func FindUser(db *sql.DB, email string) (User, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM wn_user WHERE email = '%s';", email))
	if err != nil { return User{}, err }
	users, err := readUsers(rows)
	if err != nil { return User{}, err}
	if len(users) == 0 { return User{}, misc.NotFoundError }
	return users[0], nil
}