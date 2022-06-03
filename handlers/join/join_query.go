package join

import (
	"wellnus/backend/handlers/httpError"
	"database/sql"
	"fmt"
)

// Helper function

func readJoinRequests(rows *sql.Rows) ([]JoinRequest, error) {
	joinRequests := make([]JoinRequest, 0)
	for rows.Next() {
		var joinRequest JoinRequest
		if err := rows.Scan(&joinRequest.ID, &joinRequest.UserID, &joinRequest.GroupID, &joinRequest.RequestStatus); err != nil {
			return nil, err
		}
		joinRequests = append(joinRequests, joinRequest)
	}
	return joinRequests, nil
}

func readJoinRequestWithGroups(rows *sql.Rows) ([]JoinRequestWithGroup, error) {
	joinRequestWithGroups := make([]JoinRequestWithGroup, 0)
	for rows.Next() {
		var joinRequestWithGroup JoinRequestWithGroup
		err := rows.Scan(
			&joinRequestWithGroup.JoinRequest.ID, 
			&joinRequestWithGroup.JoinRequest.UserID, 
			&joinRequestWithGroup.JoinRequest.GroupID, 
			&joinRequestWithGroup.JoinRequest.RequestStatus,
			&joinRequestWithGroup.Group.ID, 
			&joinRequestWithGroup.Group.GroupName, 
			&joinRequestWithGroup.Group.GroupDescription, 
			&joinRequestWithGroup.Group.Category, 
			&joinRequestWithGroup.Group.OwnerID)
		if err != nil {
			return nil, err
		}
		joinRequestWithGroups = append(joinRequestWithGroups, joinRequestWithGroup)
	}
	return joinRequestWithGroups, nil
}

func loadLastJoinRequestID(db *sql.DB, joinRequest JoinRequest) (JoinRequest, error) {
	row, err := db.Query("SELECT last_value FROM wn_join_request_id_seq;")
	if err != nil { return JoinRequest{}, err }
	defer row.Close()

	row.Next()
	if err := row.Scan(&joinRequest.ID); err != nil { return JoinRequest{}, err }
	return joinRequest, nil
}

func getJoinRequestWithGroup(db *sql.DB, joinRequestID int64) (JoinRequestWithGroup, error) {
	query := fmt.Sprintf(
		`SELECT * FROM wn_join_request
			JOIN wn_group
			ON wn_group.id = wn_join_request.group_id
			WHERE wn_join_request.id = %d`,
		joinRequestID)
	rows, err := db.Query(query)
	if err != nil { return JoinRequestWithGroup{}, err }
	joinRequestWithGroups, err := readJoinRequestWithGroups(rows)
	if err != nil { return JoinRequestWithGroup{}, err }
	if len(joinRequestWithGroups) == 0 { return JoinRequestWithGroup{}, httpError.NotFoundError }
	return joinRequestWithGroups[0], nil
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

// Main function

func GetAllJoinRequestsSent(db *sql.DB, userID int64) ([]JoinRequest, error) {
	query := fmt.Sprintf(
		`SELECT 
			id, 
			user_id, 
			group_id,
			request_status
		FROM wn_join_request
		WHERE user_id = %d`,
		userID)
	rows, err := db.Query(query)
	if err != nil { return nil, err }
	joinRequests, err := readJoinRequests(rows)
	if err != nil { return nil, err }
	return joinRequests, nil
}

func GetAllJoinRequestsReceived(db *sql.DB, userID int64) ([]JoinRequest, error) {
	query := fmt.Sprintf(
		`SELECT 
			wn_join_request.id, 
			wn_join_request.user_id, 
			wn_join_request.group_id,
			wn_join_request.request_status
		FROM wn_join_request
		JOIN wn_group ON wn_group.id = wn_join_request.group_id
		WHERE wn_group.owner_id = %d`,
		userID)
	rows, err := db.Query(query)
	if err != nil { return nil, err }
	joinRequests, err := readJoinRequests(rows)
	if err != nil { return nil, err }
	return joinRequests, nil
}

func GetAllJoinRequests(db *sql.DB, userID int64) ([]JoinRequest, error) {
	query := fmt.Sprintf(
		`SELECT 
			wn_join_request.id, 
			wn_join_request.user_id, 
			wn_join_request.group_id,
			wn_join_request.request_status
		FROM wn_join_request
		JOIN wn_group ON wn_group.id = wn_join_request.group_id
		WHERE wn_group.owner_id = %d OR wn_join_request.user_id = %d`,
		userID,
		userID)
	rows, err := db.Query(query)
	if err != nil { return nil, err }
	joinRequests, err := readJoinRequests(rows)
	if err != nil { return nil, err }
	return joinRequests, nil
}

func GetJoinRequest(db *sql.DB, joinRequestID int64) (JoinRequest, error) {
	query := fmt.Sprintf("SELECT * FROM wn_join_request WHERE id = %d", joinRequestID)
	rows, err := db.Query(query)
	if err != nil { return JoinRequest{}, err }
	joinRequests, err := readJoinRequests(rows)
	if err != nil { return JoinRequest{}, err }
	if len(joinRequests) == 0 { return JoinRequest{}, httpError.NotFoundError }
	return joinRequests[0], nil
}

func AddJoinRequest(db *sql.DB, groupID int64, userID int64) (JoinRequest, error) {
	query := fmt.Sprintf(
		`INSERT INTO wn_join_request (
			user_id, 
			group_id, 
			request_status
		) values (%d, %d, 'PENDING');`, 
		userID,
		groupID)
	_, err := db.Query(query)
	if err != nil { return JoinRequest{}, err }
	joinRequest, err := loadLastJoinRequestID(db, JoinRequest{ UserID: userID, GroupID: groupID, RequestStatus: "PENDING" })
	if err != nil { return JoinRequest{}, err }
	return joinRequest, nil
}

func RespondJoinRequest(db *sql.DB, joinRequestID int64, userID int64, approve bool) (JoinRequest, error) {
	joinRequestWithGroup, err := getJoinRequestWithGroup(db, joinRequestID)
	if err != nil { return JoinRequest{}, nil }
	if joinRequestWithGroup.Group.OwnerID != userID { return JoinRequest{}, httpError.UnauthorizedError }
	
	//Adding user into group if necessary
	if approve { 
		if err = addUserToGroup(db, joinRequestWithGroup.Group.ID, joinRequestWithGroup.JoinRequest.UserID); err != nil {
			return JoinRequest{}, err
		}
	}

	// Updating of status
	newStatus := "REJECTED"
	if approve {
		newStatus = "APPROVED"
	}
	joinRequestWithGroup.JoinRequest.RequestStatus = newStatus
	query := fmt.Sprintf(
		`UPDATE wn_join_request SET
			request_status = '%s'
		WHERE id = %d`,
		newStatus,
		joinRequestID)
	_, err = db.Query(query)
	if err != nil { return JoinRequest{}, err }

	return joinRequestWithGroup.JoinRequest, nil
}

func DeleteJoinRequest(db *sql.DB, joinRequestID int64, userID int64) (JoinRequest, error) {
	joinRequest, err := GetJoinRequest(db, joinRequestID)
	if err != nil { return JoinRequest{}, err }
	if joinRequest.UserID != userID { return JoinRequest{}, httpError.UnauthorizedError }

	query := fmt.Sprintf("DELETE FROM wn_join_request WHERE id = %d", joinRequestID)
	_, err = db.Query(query)
	if err != nil { return JoinRequest{}, err }
	return JoinRequest{ ID : joinRequestID }, nil
}