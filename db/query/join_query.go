package query

import (
	"wellnus/backend/handlers/misc"
	"wellnus/backend/db/model"

	"database/sql"
	"fmt"
)

type JoinRequest = model.JoinRequest
type LoadedJoinRequest = model.LoadedJoinRequest
type JoinRequestRespond = model.JoinRequestRespond

// Helper function

func readJoinRequests(rows *sql.Rows) ([]JoinRequest, error) {
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

func loadLastJoinRequestID(db *sql.DB, joinRequest JoinRequest) (JoinRequest, error) {
	row, err := db.Query("SELECT last_value FROM wn_join_request_id_seq;")
	if err != nil { return JoinRequest{}, err }
	defer row.Close()

	row.Next()
	if err := row.Scan(&joinRequest.ID); err != nil { return JoinRequest{}, err }
	return joinRequest, nil
}

func loadJoinRequest(db *sql.DB, joinRequest JoinRequest) (LoadedJoinRequest, error) {
	user, err := getUser(db, joinRequest.UserID)
	if err != nil { return LoadedJoinRequest{}, err }
	group, err := getGroup(db, joinRequest.GroupID)
	if err != nil { return LoadedJoinRequest{}, err }
	return LoadedJoinRequest{ JoinRequest: joinRequest, User: user, Group: group }, nil
}

func getJoinRequest(db *sql.DB, joinRequestID int64) (JoinRequest, error) {
	query := fmt.Sprintf("SELECT * FROM wn_join_request WHERE id = %d", joinRequestID)
	rows, err := db.Query(query)
	if err != nil { return JoinRequest{}, err }
	joinRequests, err := readJoinRequests(rows)
	if err != nil { return JoinRequest{}, err }
	if len(joinRequests) == 0 { return JoinRequest{}, misc.NotFoundError }
	return joinRequests[0], nil
}

// Main function

func GetAllJoinRequestsSentOfUser(db *sql.DB, userID int64) ([]JoinRequest, error) {
	query := fmt.Sprintf(
		`SELECT 
			id, 
			user_id, 
			group_id
		FROM wn_join_request
		WHERE user_id = %d`,
		userID)
	rows, err := db.Query(query)
	if err != nil { return nil, err }
	joinRequests, err := readJoinRequests(rows)
	if err != nil { return nil, err }
	return joinRequests, nil
}

func GetAllJoinRequestsReceivedOfUser(db *sql.DB, userID int64) ([]JoinRequest, error) {
	query := fmt.Sprintf(
		`SELECT 
			wn_join_request.id, 
			wn_join_request.user_id, 
			wn_join_request.group_id
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

func GetAllJoinRequestsOfUser(db *sql.DB, userID int64) ([]JoinRequest, error) {
	query := fmt.Sprintf(
		`SELECT 
			wn_join_request.id, 
			wn_join_request.user_id, 
			wn_join_request.group_id
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

func GetLoadedJoinRequest(db *sql.DB, joinRequestID int64) (LoadedJoinRequest, error) {
	joinRequest, err := getJoinRequest(db, joinRequestID)
	if err != nil { return LoadedJoinRequest{}, err }
	loadedJoinRequest, err := loadJoinRequest(db, joinRequest)
	if err != nil { return LoadedJoinRequest{}, err }
	return loadedJoinRequest, nil
}

func AddJoinRequest(db *sql.DB, groupID int64, userID int64) (JoinRequest, error) {
	query := fmt.Sprintf(
		`INSERT INTO wn_join_request (
			user_id, 
			group_id
		) values (%d, %d);`, 
		userID,
		groupID)
	_, err := db.Query(query)
	if err != nil { return JoinRequest{}, err }
	joinRequest, err := loadLastJoinRequestID(db, JoinRequest{ UserID: userID, GroupID: groupID })
	if err != nil { return JoinRequest{}, err }
	return joinRequest, nil
}

func RespondJoinRequest(db *sql.DB, joinRequestID int64, userID int64, approve bool) (JoinRequestRespond, error) {
	loadedJoinRequest, err := GetLoadedJoinRequest(db, joinRequestID)
	if err != nil { return JoinRequestRespond{}, nil }
	if loadedJoinRequest.Group.OwnerID != userID { return JoinRequestRespond{}, misc.UnauthorizedError }
	
	//Adding user into group if necessary
	if approve { 
		if err = addUserToGroup(db, loadedJoinRequest.Group.ID, loadedJoinRequest.JoinRequest.UserID); err != nil {
			return JoinRequestRespond{}, err
		}
	}
	query := fmt.Sprintf("DELETE FROM wn_join_request WHERE id = %d", joinRequestID)
	_, err = db.Query(query)
	if err != nil { return JoinRequestRespond{}, err }
	return JoinRequestRespond{ Approve: approve }, nil
}

func DeleteJoinRequest(db *sql.DB, joinRequestID int64, userID int64) (JoinRequest, error) {
	joinRequest, err := getJoinRequest(db, joinRequestID)
	if err != nil { return JoinRequest{}, err }
	if joinRequest.UserID != userID { return JoinRequest{}, misc.UnauthorizedError }

	query := fmt.Sprintf("DELETE FROM wn_join_request WHERE id = %d", joinRequestID)
	_, err = db.Query(query)
	if err != nil { return JoinRequest{}, err }
	return JoinRequest{ ID : joinRequestID }, nil
}