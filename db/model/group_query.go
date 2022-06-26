package model

import (
	"wellnus/backend/router/http_helper/http_error"	
	"log"
	"fmt"
	"database/sql"
	"errors"
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

func GetGroup(db *sql.DB, groupID int64) (Group, error) {
	rows, err := db.Query("SELECT * FROM wn_group WHERE id = $1;", groupID)
	if err != nil { return Group{}, err }
	defer rows.Close()
	groups, err := ReadGroups(rows)
	if err != nil { return Group{}, err }
	if len(groups) == 0 { return Group{}, http_error.NotFoundError }
	return groups[0], nil
}

func ChangeOwnership(db *sql.DB, group Group, newOwnerID int64) (Group, error) {
	group.OwnerID = newOwnerID
	_, err := db.Exec(
		`UPDATE wn_group SET 
			owner_id = $1
		WHERE id = $2;`,
		group.OwnerID,
		group.ID)
	if err != nil { return Group{}, err }
	return group, nil
}

func AddUserToGroup(db *sql.DB, groupID int64, userID int64) error {
	_, err := db.Exec(
		`INSERT INTO wn_user_group (
			user_id, 
			group_id) 
		VALUES ($1, $2)`, 
		userID, 
		groupID)
	return err
}

func RemoveUserFromGroup(db *sql.DB, groupID int64, userID int64) error {
	_, err := db.Exec(
		`DELETE FROM wn_user_group WHERE
			user_id = $1 AND
			group_id = $2`,
		userID,
		groupID)
	return err
}

func deleteGroup(db *sql.DB, groupID int64) error {
	_, err := db.Exec("DELETE FROM wn_group WHERE id = $1", groupID)
	return err
}

// Main Functions

func GetGroupWithUsers(db *sql.DB, groupID int64) (GroupWithUsers, error) {
	group, err := GetGroup(db, groupID)
	if err != nil { return GroupWithUsers{}, err }
	users, err := GetAllUsersOfGroup(db, groupID)
	if err != nil { return GroupWithUsers{}, err }
	return GroupWithUsers{ Group: group, Users: users}, nil
}

func GetAllGroupsOfUser(db *sql.DB, userID int64) ([]Group, error) {
	rows, err := db.Query(
		`SELECT
			wn_group.id, 
			wn_group.group_name, 
			wn_group.group_description,
			wn_group.category, 
			wn_group.owner_id
		FROM wn_user_group JOIN wn_group 
		ON wn_user_group.group_id = wn_group.id 
		WHERE wn_user_group.user_id = $1`,
		userID)
	if err != nil { return nil, err }
	defer rows.Close()
	groups, err := ReadGroups(rows)
	if err != nil { return nil, err}
	return groups, nil
}

func AddGroupWithUserIDs(db *sql.DB, group Group, userIDs []int64) (GroupWithUsers, error) {
	if len(userIDs) == 0 { return GroupWithUsers{}, errors.New("Insufficient users to form a group") }
	ownerID := userIDs[0]
	userIDs = userIDs[1:]
	group.OwnerID = ownerID //Taking first userID as ownerID
	_, err := db.Exec(
		`INSERT INTO wn_group (
			group_name, 
			group_description, 
			category, 
			owner_id) 
		VALUES ($1, $2, $3, $4);`,
		group.GroupName,
		group.GroupDescription,
		group.Category,
		group.OwnerID)
	if err != nil { return GroupWithUsers{}, err }
	group, err = group.LoadLastGroupID(db)
	if err != nil { return GroupWithUsers{}, err }

	// Adding Owner
	if err = AddUserToGroup(db, group.ID, ownerID); err != nil {
		log.Printf("Failed to add Owner: %v", err)
		if _, fatal := db.Exec("DELETE FROM wn_group WHERE id = $1", group.ID); fatal != nil {
			log.Fatal(fmt.Sprintf("Failed to remove added group after failing to add owner. Fatal: %v", fatal))
		}
		return GroupWithUsers{}, err
	}

	// Adding Other Users
	for _, userID := range userIDs {
		if err := AddUserToGroup(db, group.ID, userID); err != nil {
			return GroupWithUsers{}, err
		}
	}
	users, err := GetAllUsersOfGroup(db, group.ID)
	if err != nil { return GroupWithUsers{}, err }

	return GroupWithUsers{ Group: group, Users: users }, nil
}

func UpdateGroup(db *sql.DB, updatedGroup Group, groupID int64, userID int64) (Group, error) {
	targetGroup, err := GetGroup(db, groupID)
	if err != nil { return Group{}, err }
	if targetGroup.OwnerID != userID { return Group{}, http_error.UnauthorizedError }

	updatedGroup = updatedGroup.MergeGroup(targetGroup)

	_, err = db.Exec(
		`UPDATE wn_group SET 
			group_name = $1,
			group_description = $2,
			category = $3,
			owner_id = $4
		WHERE id = $5;`,
		updatedGroup.GroupName,
		updatedGroup.GroupDescription,
		updatedGroup.Category,
		updatedGroup.OwnerID,
		groupID)
	if err != nil { return Group{}, err }
	return updatedGroup, nil
}

func LeaveGroup(db *sql.DB, groupID int64, userID int64) (GroupWithUsers, error) {
	targetGroupWithUsers, err := GetGroupWithUsers(db, groupID)
	if err != nil { return GroupWithUsers{}, err }
	if targetGroupWithUsers.Group.OwnerID == userID {
		newOwnerID := targetGroupWithUsers.GetNewOwnerID()
		if newOwnerID == 0 {	
			err = deleteGroup(db, groupID)
			if err != nil { return GroupWithUsers{}, err } // Group not deleted
			return GroupWithUsers{ Group: Group{ID: groupID} }, nil
		}
		targetGroupWithUsers.Group, err = ChangeOwnership(db, targetGroupWithUsers.Group, newOwnerID)
		if err != nil { return GroupWithUsers{}, err } // Ownership not transferred
	}
	err = RemoveUserFromGroup(db, groupID, userID)
	if err != nil { return GroupWithUsers{}, err } // User not properly removed
	users, err := GetAllUsersOfGroup(db, groupID)
	if err != nil { return GroupWithUsers{}, err }
	targetGroupWithUsers.Users = users
	if err != nil { return GroupWithUsers{}, err } // reloading of group with users failed
	return targetGroupWithUsers, nil
}

func LeaveAllGroups(db *sql.DB, userID int64) ([]GroupWithUsers, error) {
	groups, err := GetAllGroupsOfUser(db, userID)
	if err != nil { return nil, err}
	groupsWithUsers := make([]GroupWithUsers, 0)
	for _, group := range groups {
		groupWithUsers, err := LeaveGroup(db, group.ID, userID)
		if err != nil { return nil, err}
		groupsWithUsers = append(groupsWithUsers, groupWithUsers)
	}
	return groupsWithUsers, nil
}

func IsUserInGroup(db *sql.DB, userID int64, groupID int64) (bool, error) {
	row, err := db.Query(
		`SELECT COUNT(*) != 0 FROM wn_user_group 
		WHERE user_id = $1 and group_id = $2`,
		userID,
		groupID)
	if err != nil { return false, err }
	defer row.Close()
	var membership bool
	row.Next()
	if err := row.Scan(&membership); err != nil { return false, nil }
	return membership, nil
}