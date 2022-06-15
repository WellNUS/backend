package model

import (
	"database/sql"
)

type JoinRequestRespond struct {
	Approve bool `json:"approve"`
}

type JoinRequest struct {
	ID 				int64 	`json:"id"`
	UserID 			int64 	`json:"user_id"`
	GroupID 		int64 	`json:"group_id"`
}

type LoadedJoinRequest struct {
	JoinRequest		JoinRequest 	`json:"join_request"`
	User			User			`json:"user"`
	Group			Group			`json:"group"`
}

func (joinRequest JoinRequest) LoadJoinRequest(db *sql.DB) (LoadedJoinRequest, error) {
	user, err := GetUser(db, joinRequest.UserID)
	if err != nil { return LoadedJoinRequest{}, err }
	group, err := GetGroup(db, joinRequest.GroupID)
	if err != nil { return LoadedJoinRequest{}, err }
	return LoadedJoinRequest{ JoinRequest: joinRequest, User: user, Group: group }, nil
}

func (joinRequest JoinRequest) LoadLastJoinRequestID(db *sql.DB) (JoinRequest, error) {
	row, err := db.Query("SELECT last_value FROM wn_join_request_id_seq;")
	if err != nil { return JoinRequest{}, err }
	defer row.Close()

	row.Next()
	if err := row.Scan(&joinRequest.ID); err != nil { return JoinRequest{}, err }
	return joinRequest, nil
}

func (joinRequest1 JoinRequest) Equal(joinRequest2 JoinRequest) bool {
	return joinRequest1.ID == joinRequest2.ID &&
	joinRequest1.GroupID == joinRequest2.GroupID &&
	joinRequest1.UserID == joinRequest2.UserID
}