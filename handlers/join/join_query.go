package join

import (
	"wellnus/backend/references"
	"database/sql"
)

type JoinRequest = references.JoinRequest

func LoadLastJoinRequestID(db *sql.DB, joinRequest JoinRequest) (JoinRequest, error) {
	row, err := db.Query("SELECT last_value FROM wn_join_request_id_seq;")
	if err != nil { return User{}, err }
	defer row.Close()

	row.Next()
	if err := row.Scan(&joinRequest.ID); err != nil { return JoinRequest{}, err }
	return joinRequest, nil
}

func AddJoinRequest(db *sql.DB, groupID int64, userID int64) (JoinRequest, error) {
	_, err := db.Query(fmt.Sprintf(
		"INSERT INTO wn_join_request (user_id, group_id, request_status) values (%d, %d, 'PENDING');", 
		userID,
		groupID))
	if err != nil { return JoinRequest{}, err }
	joinRequest, err := JoinRequest(db, JoinRequest{ UserID: userID, GroupID: groupID, RequestStatus: "PENDING" })
	if err != nil { return JoinRequest{}, err }
	return joinRequest, nil
}