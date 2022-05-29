package group;

import (
	"wellnus/backend/handlers/httpError"
	
	"log"
	"fmt"
	"database/sql"
)

// Helper function
func ReadGroups(rows *sql.Rows) ([]Group, error) {
	groups := make([]Group, 0)
	for rows.Next() {
		var group Group
		if err := rows.Scan(&group.ID, &group.GroupName, &group.GroupDescription, &group.Category, &group.OwnerID); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func MergeGroup(groupMain Group, groupAdd Group) Group {
	groupMain.ID = groupAdd.ID
	if groupMain.GroupName == "" {
		groupMain.GroupName = groupAdd.GroupName
	}
	if groupMain.GroupDescription == "" {
		groupMain.GroupDescription = groupAdd.GroupDescription
	}
	if groupMain.Category == "" {
		groupMain.Category = groupAdd.Category
	}
	if groupMain.OwnerID == 0 {
		groupMain.OwnerID = groupAdd.OwnerID
	}
	return groupMain
}

func GetUsersInGroup(db *sql.DB, id int64) ([]User, error) {
	rows, err := db.Query(fmt.Sprintf(
		"SELECT * FROM wn_user_group JOIN wn_user ON wn_user_group.user_id = wn_user.id WHERE wn_user_group.group_id = %d", 
		id))
	if err != nil { return nil, err}
	users := make([]User, 0)
	for rows.Next() {
		var userID, groupID int64; // Temp variables
		var user User
		if err := rows.Scan(&userID, &groupID, &user.ID, &user.FirstName, &user.LastName, &user.Gender, &user.Email, &user.UserRole, &user.PasswordHash); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Main Functions
func GetGroup(db *sql.DB, id int64) (GroupWithUsers, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM wn_group WHERE id = %d;", id))
	if err != nil { return GroupWithUsers{}, err }
	defer rows.Close()

	groups, err := ReadGroups(rows)
	if err != nil { return GroupWithUsers{}, err }
	if len(groups) == 0 { return GroupWithUsers{}, httpError.NotFoundError }
	group := groups[0]

	users, err := GetUsersInGroup(db, id)
	if err != nil { return GroupWithUsers{}, err }
	return GroupWithUsers{ Group: group, Users: users }, nil
}

func GetAllGroups(db *sql.DB) ([]Group, error) {
	rows, err := db.Query("SELECT * FROM wn_group;")
	if err != nil { return nil, err }
	defer rows.Close()
	
	groups, err := ReadGroups(rows)
	if err != nil { return nil, err}
	return groups, nil
}

func AddGroup(db *sql.DB, newGroup Group) (Group, error) {
	_, err := db.Query(fmt.Sprintf(
		"INSERT INTO wn_group (group_name, group_description, category, owner_id) VALUES ('%s', '%s', '%s', %d);",
		newGroup.GroupName,
		newGroup.GroupDescription,
		newGroup.Category,
		newGroup.OwnerID))
	if err != nil { return Group{}, err }
	row, err := db.Query("SELECT last_value FROM wn_group_id_seq;")
	if err != nil { return Group{}, err }
	defer row.Close()
	row.Next()
	if err := row.Scan(&newGroup.ID); err != nil { return Group{}, err }

	// new Group successfull added into DB. Now adding owner into new group
	err = AddUserToGroup(db, newGroup.OwnerID, newGroup.ID)
	if err != nil {
		log.Printf("Failed to add Owner: %v", err)
		if _, fatal := db.Query(fmt.Sprintf("DELETE FROM wn_group WHERE id = %d", newGroup.ID)); fatal != nil {
			log.Fatal(fmt.Sprintf("Failed to remove added group after failing to add owner. Fatal: %v", fatal))
		}
		return Group{}, err
	}
	return newGroup, nil
}

func AddUserToGroup(db *sql.DB, userID int64, groupID int64) error {
	_, err := db.Query(fmt.Sprintf(
		"INSERT INTO wn_user_group (user_id, group_id) VALUES (%d, %d)", 
		userID, 
		groupID))
	return err
}

/*
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
*/