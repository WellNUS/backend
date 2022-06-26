package model

import (
	"database/sql"
)

type Group struct {
	ID					int64	`json:"id"`
	GroupName			string	`json:"group_name"`
	GroupDescription 	string	`json:"group_description"`
	Category			string 	`json:"category"`
	OwnerID				int64	`json:"owner_id"`
}

type GroupWithUsers struct {
	Group	Group	`json:"group"`
	Users	[]User	`json:"users"`
}

func (groupMain Group) MergeGroup(groupAdd Group) Group {
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

func (group Group) LoadLastGroupID(db *sql.DB) (Group, error) {
	row, err := db.Query("SELECT last_value FROM wn_group_id_seq;")
	defer row.Close()
	if err != nil { return Group{}, err }

	row.Next()
	if err := row.Scan(&group.ID); err != nil { return Group{}, err }
	return group, nil
}

func (groupWithUsers GroupWithUsers) GetNewOwnerID() int64 {
	currOwnerID := groupWithUsers.Group.OwnerID
	users := groupWithUsers.Users
	for _, user := range users {
		if user.ID != currOwnerID {
			return user.ID
		}
	}
	return 0
}

func (group1 Group) Equal(group2 Group) bool {
	return group1.ID == group2.ID &&
	group1.GroupName == group2.GroupName &&
	group1.GroupDescription == group2.GroupDescription &&
	group1.Category == group2.Category &&
	group1.OwnerID == group2.OwnerID
}