package model

import (
	"wellnus/backend/router/misc/http_error"
	"database/sql"
)

// Helper function

func ReadJoinRequests(rows *sql.Rows) ([]JoinRequest, error) {
	joinRequests := make([]JoinRequest, 0)
	for rows.Next() {
		var joinRequest JoinRequest
		if err := rows.Scan(&joinRequest.ID, &joinRequest.UserID, &joinRequest.GroupID); err != nil {
			return nil, err
		}
		joinRequests = append(joinRequests, joinRequest)
	}
	return joinRequests, nil
}

func GetJoinRequest(db *sql.DB, joinRequestID int64) (JoinRequest, error) {
	rows, err := db.Query("SELECT * FROM wn_join_request WHERE id = $1", joinRequestID)
	if err != nil { return JoinRequest{}, err }
	joinRequests, err := ReadJoinRequests(rows)
	if err != nil { return JoinRequest{}, err }
	if len(joinRequests) == 0 { return JoinRequest{}, http_error.NotFoundError }
	return joinRequests[0], nil
}

// Main function

func GetAllJoinRequestsSentOfUser(db *sql.DB, userID int64) ([]JoinRequest, error) {
	rows, err := db.Query(
		`SELECT 
			id, 
			user_id, 
			group_id
		FROM wn_join_request
		WHERE user_id = $1`,
		userID)
	if err != nil { return nil, err }
	joinRequests, err := ReadJoinRequests(rows)
	if err != nil { return nil, err }
	return joinRequests, nil
}

func GetAllJoinRequestsReceivedOfUser(db *sql.DB, userID int64) ([]JoinRequest, error) {
	rows, err := db.Query(
		`SELECT 
			wn_join_request.id, 
			wn_join_request.user_id, 
			wn_join_request.group_id
		FROM wn_join_request
		JOIN wn_group ON wn_group.id = wn_join_request.group_id
		WHERE wn_group.owner_id = $1`,
		userID)
	if err != nil { return nil, err }
	joinRequests, err := ReadJoinRequests(rows)
	if err != nil { return nil, err }
	return joinRequests, nil
}

func GetAllJoinRequestsOfUser(db *sql.DB, userID int64) ([]JoinRequest, error) {
	rows, err := db.Query(
		`SELECT 
			wn_join_request.id, 
			wn_join_request.user_id, 
			wn_join_request.group_id
		FROM wn_join_request
		JOIN wn_group ON wn_group.id = wn_join_request.group_id
		WHERE wn_group.owner_id = $1 OR wn_join_request.user_id = $2`,
		userID,
		userID)
	if err != nil { return nil, err }
	joinRequests, err := ReadJoinRequests(rows)
	if err != nil { return nil, err }
	return joinRequests, nil
}

func GetLoadedJoinRequest(db *sql.DB, joinRequestID int64) (LoadedJoinRequest, error) {
	joinRequest, err := GetJoinRequest(db, joinRequestID)
	if err != nil { return LoadedJoinRequest{}, err }
	loadedJoinRequest, err := joinRequest.LoadJoinRequest(db)
	if err != nil { return LoadedJoinRequest{}, err }
	return loadedJoinRequest, nil
}

func AddJoinRequest(db *sql.DB, groupID int64, userID int64) (JoinRequest, error) {
	_, err := db.Exec(
		`INSERT INTO wn_join_request (
			user_id, 
			group_id
		) values ($1, $2);`, 
		userID,
		groupID)
	if err != nil { return JoinRequest{}, err }
	joinRequest, err := JoinRequest{ UserID: userID, GroupID: groupID }.LoadLastJoinRequestID(db)
	if err != nil { return JoinRequest{}, err }
	return joinRequest, nil
}

func RespondJoinRequest(db *sql.DB, joinRequestID int64, userID int64, approve bool) (JoinRequestRespond, error) {
	loadedJoinRequest, err := GetLoadedJoinRequest(db, joinRequestID)
	if err != nil { return JoinRequestRespond{}, nil }
	if loadedJoinRequest.Group.OwnerID != userID { return JoinRequestRespond{}, http_error.UnauthorizedError }
	
	//Adding user into group if necessary
	if approve { 
		if err = AddUserToGroup(db, loadedJoinRequest.Group.ID, loadedJoinRequest.JoinRequest.UserID); err != nil {
			return JoinRequestRespond{}, err
		}
	}
	_, err = db.Exec("DELETE FROM wn_join_request WHERE id = $1", joinRequestID)
	if err != nil { return JoinRequestRespond{}, err }
	return JoinRequestRespond{ Approve: approve }, nil
}

func DeleteJoinRequest(db *sql.DB, joinRequestID int64, userID int64) (JoinRequest, error) {
	joinRequest, err := GetJoinRequest(db, joinRequestID)
	if err != nil { return JoinRequest{}, err }
	if joinRequest.UserID != userID { return JoinRequest{}, http_error.UnauthorizedError }

	_, err = db.Exec("DELETE FROM wn_join_request WHERE id = $1", joinRequestID)
	if err != nil { return JoinRequest{}, err }
	return JoinRequest{ ID : joinRequestID }, nil
}