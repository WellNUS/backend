package query

import (
	"wellnus/backend/handlers/misc"
	"wellnus/backend/db/model"
	
	"log"
	"fmt"
	"database/sql"
)

type Group = model.Group
type GroupWithUsers = model.GroupWithUsers

// Helper function
func mergeGroup(groupMain Group, groupAdd Group) Group {
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

func readGroups(rows *sql.Rows) ([]Group, error) {
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

func getGroup(db *sql.DB, groupID int64) (Group, error) {
	query := fmt.Sprintf("SELECT * FROM wn_group WHERE id = %d;", groupID)
	rows, err := db.Query(query)
	if err != nil { return Group{}, err }
	defer rows.Close()

	groups, err := readGroups(rows)
	if err != nil { return Group{}, err }
	if len(groups) == 0 { return Group{}, misc.NotFoundError }
	return groups[0], nil
}

func loadLastGroupID(db *sql.DB, group Group) (Group, error) {
	row, err := db.Query("SELECT last_value FROM wn_group_id_seq;")
	if err != nil { return Group{}, err }
	defer row.Close()

	row.Next()
	if err := row.Scan(&group.ID); err != nil { return Group{}, err }
	return group, nil
}

func changeOwnership(db *sql.DB, group Group, newOwnerID int64) (Group, error) {
	group.OwnerID = newOwnerID
	query := fmt.Sprintf(
		`UPDATE wn_group SET 
			owner_id = %d
		WHERE id = %d;`,
		group.OwnerID,
		group.ID)
	_, err := db.Query(query)
	if err != nil { return Group{}, err }
	return group, nil
}

func addUserToGroup(db *sql.DB, groupID int64, userID int64) error {
	query := fmt.Sprintf(
		`INSERT INTO wn_user_group (
			user_id, 
			group_id) 
		VALUES (%d, %d)`, 
		userID, 
		groupID)
	_, err := db.Query(query)
	return err
}

func removeUserFromGroup(db *sql.DB, groupID int64, userID int64) error {
	query := fmt.Sprintf(
		`DELETE FROM wn_user_group WHERE
			user_id = %d AND
			group_id = %d`,
			userID,
			groupID)
	_, err := db.Query(query)
	return err
} 

func getNewOwnerID(groupWithUsers GroupWithUsers) int64 {
	currOwnerID := groupWithUsers.Group.OwnerID
	users := groupWithUsers.Users
	for _, user := range users {
		if user.ID != currOwnerID {
			return user.ID
		}
	}
	return 0
}

func deleteGroup(db *sql.DB, groupID int64) error {
	query := fmt.Sprintf("DELETE FROM wn_group WHERE id = %d", groupID)
	_, err := db.Query(query)
	return err
}

// Main Functions

func GetGroupWithUsers(db *sql.DB, groupID int64) (GroupWithUsers, error) {
	group, err := getGroup(db, groupID)
	if err != nil { return GroupWithUsers{}, err }
	users, err := GetAllUsersOfGroup(db, groupID)
	if err != nil { return GroupWithUsers{}, err }
	return GroupWithUsers{ Group: group, Users: users}, nil
}

func GetAllGroupsOfUser(db *sql.DB, userID int64) ([]Group, error) {
	query := fmt.Sprintf(
		`SELECT
			wn_group.id, 
			wn_group.group_name, 
			wn_group.group_description,
			wn_group.category, 
			wn_group.owner_id
		FROM wn_user_group JOIN wn_group 
		ON wn_user_group.group_id = wn_group.id 
		WHERE wn_user_group.user_id = %d`,
		userID)
	rows, err := db.Query(query)
	if err != nil { return nil, err }
	defer rows.Close()
	
	groups, err := readGroups(rows)
	if err != nil { return nil, err}
	return groups, nil
}

func AddGroup(db *sql.DB, newGroup Group) (GroupWithUsers, error) {
	query := fmt.Sprintf(
		`INSERT INTO wn_group (
			group_name, 
			group_description, 
			category, 
			owner_id) 
		VALUES ('%s', '%s', '%s', %d);`,
		newGroup.GroupName,
		newGroup.GroupDescription,
		newGroup.Category,
		newGroup.OwnerID)
	_, err := db.Query(query)
	if err != nil { return GroupWithUsers{}, err }
	newGroup, err = loadLastGroupID(db, newGroup)
	if err != nil { return GroupWithUsers{}, err }
	
	// newGroup successfully added into DB. Now adding owner into new group
	err = addUserToGroup(db, newGroup.ID, newGroup.OwnerID)
	if err != nil {
		log.Printf("Failed to add Owner: %v", err)
		if _, fatal := db.Query(fmt.Sprintf("DELETE FROM wn_group WHERE id = %d", newGroup.ID)); fatal != nil {
			log.Fatal(fmt.Sprintf("Failed to remove added group after failing to add owner. Fatal: %v", fatal))
		}
		return GroupWithUsers{}, err
	}
	users, err := GetAllUsersOfGroup(db, newGroup.ID)
	if err != nil { return GroupWithUsers{}, err }
	groupWithUsers := GroupWithUsers{ Group: newGroup, Users: users }
	if err != nil { return GroupWithUsers{}, err }
	return groupWithUsers, nil
}

func UpdateGroup(db *sql.DB, updatedGroup Group, groupID int64, userID int64) (Group, error) {
	targetGroup, err := getGroup(db, groupID)
	if err != nil { return Group{}, err }
	if targetGroup.OwnerID != userID { return Group{}, misc.UnauthorizedError }

	updatedGroup = mergeGroup(updatedGroup, targetGroup)

	query := fmt.Sprintf(
		`UPDATE wn_group SET 
			group_name = '%s',
			group_description = '%s',
			category = '%s',
			owner_id = %d
		WHERE id = %d;`,
		updatedGroup.GroupName,
		updatedGroup.GroupDescription,
		updatedGroup.Category,
		updatedGroup.OwnerID,
		groupID)
	if _, err := db.Query(query); err != nil {
		return Group{}, err;
	}
	return updatedGroup, nil;
}

func LeaveGroup(db *sql.DB, groupID int64, userID int64) (GroupWithUsers, error) {
	targetGroupWithUsers, err := GetGroupWithUsers(db, groupID)
	if err != nil { return GroupWithUsers{}, err }
	if targetGroupWithUsers.Group.OwnerID == userID {
		newOwnerID := getNewOwnerID(targetGroupWithUsers)
		if newOwnerID == 0 {	
			err = deleteGroup(db, groupID)
			if err != nil { return GroupWithUsers{}, err } // Group not deleted
			return GroupWithUsers{ Group: Group{ID: groupID} }, nil
		}
		targetGroupWithUsers.Group, err = changeOwnership(db, targetGroupWithUsers.Group, newOwnerID)
		if err != nil { return GroupWithUsers{}, err } // Ownership not transferred
	}
	err = removeUserFromGroup(db, groupID, userID)
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