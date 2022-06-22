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

func ReadLoadedJoinRequests(rows *sql.Rows) ([]LoadedJoinRequest, error) {
	loadedJoinRequests := make([]LoadedJoinRequest, 0)
	for rows.Next() {
		var loadedJoinRequest LoadedJoinRequest
		if err := rows.Scan(
			&loadedJoinRequest.JoinRequest.ID,
			&loadedJoinRequest.JoinRequest.UserID,
			&loadedJoinRequest.JoinRequest.GroupID,
			&loadedJoinRequest.User.ID,
			&loadedJoinRequest.User.FirstName,
			&loadedJoinRequest.User.LastName,
			&loadedJoinRequest.User.Gender,
			&loadedJoinRequest.User.Faculty,
			&loadedJoinRequest.User.Email,
			&loadedJoinRequest.User.UserRole,
			&loadedJoinRequest.User.PasswordHash,
			&loadedJoinRequest.Group.ID,
			&loadedJoinRequest.Group.GroupName,
			&loadedJoinRequest.Group.GroupDescription,
			&loadedJoinRequest.Group.Category,
			&loadedJoinRequest.Group.OwnerID);
			err != nil {
				return nil, err
			}
		loadedJoinRequests = append(loadedJoinRequests, loadedJoinRequest)
	}
	return loadedJoinRequests, nil
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

func GetAllLoadedJoinRequestsSentOfUser(db *sql.DB, userID int64) ([]LoadedJoinRequest, error) {
	rows, err := db.Query(
		`SELECT 
			wn_join_request.id, 
			wn_join_request.user_id, 
			wn_join_request.group_id,
			wn_user.id,
			wn_user.first_name, 
			wn_user.last_name, 
			wn_user.gender, 
			wn_user.faculty, 
			wn_user.email, 
			wn_user.user_role, 
			wn_user.password_hash,
			wn_group.id,
			wn_group.group_name, 
			wn_group.group_description, 
			wn_group.category, 
			wn_group.owner_id
		FROM wn_join_request 
		JOIN wn_user ON wn_user.id = wn_join_request.user_id
		JOIN wn_group ON wn_group.id = wn_join_request.group_id
		WHERE wn_join_request.user_id = $1`,
		userID)
	if err != nil { return nil, err }
	loadedJoinRequests, err := ReadLoadedJoinRequests(rows)
	if err != nil { return nil, err }
	return loadedJoinRequests, nil
}

func GetAllLoadedJoinRequestsReceivedOfUser(db *sql.DB, userID int64) ([]LoadedJoinRequest, error) {
	rows, err := db.Query(
		`SELECT 
			wn_join_request.id, 
			wn_join_request.user_id, 
			wn_join_request.group_id,
			wn_user.id,
			wn_user.first_name, 
			wn_user.last_name, 
			wn_user.gender, 
			wn_user.faculty, 
			wn_user.email, 
			wn_user.user_role, 
			wn_user.password_hash,
			wn_group.id,
			wn_group.group_name, 
			wn_group.group_description, 
			wn_group.category, 
			wn_group.owner_id
		FROM wn_join_request 
		JOIN wn_user ON wn_user.id = wn_join_request.user_id
		JOIN wn_group ON wn_group.id = wn_join_request.group_id
		WHERE wn_group.owner_id = $1`,
		userID)
	if err != nil { return nil, err }
	loadedJoinRequests, err := ReadLoadedJoinRequests(rows)
	if err != nil { return nil, err }
	return loadedJoinRequests, nil
}

func GetAllLoadedJoinRequestsOfUser(db *sql.DB, userID int64) ([]LoadedJoinRequest, error) {
	rows, err := db.Query(
		`SELECT 
			wn_join_request.id, 
			wn_join_request.user_id, 
			wn_join_request.group_id,
			wn_user.id,
			wn_user.first_name, 
			wn_user.last_name, 
			wn_user.gender, 
			wn_user.faculty, 
			wn_user.email, 
			wn_user.user_role, 
			wn_user.password_hash,
			wn_group.id,
			wn_group.group_name, 
			wn_group.group_description, 
			wn_group.category, 
			wn_group.owner_id
		FROM wn_join_request 
		JOIN wn_user ON wn_user.id = wn_join_request.user_id
		JOIN wn_group ON wn_group.id = wn_join_request.group_id
		WHERE wn_group.owner_id = $1 OR wn_join_request.user_id = $2`,
		userID,
		userID)
	if err != nil { return nil, err }
	loadedJoinRequests, err := ReadLoadedJoinRequests(rows)
	if err != nil { return nil, err }
	return loadedJoinRequests, nil
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