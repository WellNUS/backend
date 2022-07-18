package model

import (
	"wellnus/backend/router/http_helper/http_error"

	"errors"
	"log"
	"fmt"
	"database/sql"

	"github.com/lib/pq"
)

func ReadCounselRequests(rows *sql.Rows) ([]CounselRequest, error) {
	counselRequests := make([]CounselRequest, 0)
	for rows.Next() {
		var counselRequest CounselRequest
		if err := rows.Scan(
			&counselRequest.UserID, 
			&counselRequest.Details, 
			pq.Array(&counselRequest.Topics),
			&counselRequest.TimeAdded);
			err != nil {
				return nil, err
			}
		counselRequests = append(counselRequests, counselRequest)
	}
	return counselRequests, nil
}

func CheckCounselRequest(db *sql.DB, userID int64) (bool, error) {
	row, err := db.Query("SELECT COUNT(*) != 0 FROM wn_counsel_request WHERE user_id = $1", userID)
	if err != nil { return false, err }
	defer row.Close()
	row.Next()
	var present bool
	if err := row.Scan(&present); err != nil {
		return false, err
	}
	return present, nil
}

func AuthoriseProvider(db *sql.DB, userID int64) bool {
	user, _ := GetUser(db, userID)
	return IsProvider(user)
}

// Main functions

func GetAllCounselRequests(db *sql.DB, topics []string, userID int64) ([]CounselRequest, error) {
	authorized := AuthoriseProvider(db, userID)
	if !authorized { return nil, http_error.UnauthorizedError }
	var rows *sql.Rows
	var err error
	if topics == nil {
		rows, err = db.Query("SELECT * FROM wn_counsel_request")
	} else {
		rows, err = db.Query("SELECT * FROM wn_counsel_request WHERE $1 <@ topics", pq.Array(topics))
	}
	if err != nil { return nil, err }
	defer rows.Close()
	counselRequests, err := ReadCounselRequests(rows)
	if err != nil { return nil, err }
	return counselRequests, nil
}

func GetCounselRequest(db *sql.DB, recipientID int64, userID int64) (CounselRequest, error) {
	authorized := recipientID == userID || AuthoriseProvider(db, userID)
	if !authorized { return CounselRequest{}, http_error.UnauthorizedError }
	rows, err := db.Query("SELECT * FROM wn_counsel_request WHERE user_id = $1", recipientID)
	if err != nil { return CounselRequest{}, err }
	defer rows.Close()
	counselRequests, err := ReadCounselRequests(rows)
	if err != nil { return CounselRequest{}, err }
	if len(counselRequests) == 0 { return CounselRequest{}, http_error.NotFoundError }
	return counselRequests[0], nil
}

func AddUpdateCounselRequest(db *sql.DB, counselRequest CounselRequest, userID int64) (CounselRequest, error) {
	if counselRequest.TimeAdded.IsZero() { return CounselRequest{}, errors.New("CounselRequest had a default value for time_added") }
	counselRequest.UserID = userID
	_, err := db.Exec(
		`INSERT INTO wn_counsel_request (
			user_id,
			details,
			topics,
			time_added
		) VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			details = EXCLUDED.details,
			topics = EXCLUDED.topics,
			time_added = EXCLUDED.time_added`,
		counselRequest.UserID,
		counselRequest.Details,
		pq.Array(counselRequest.Topics),
		counselRequest.TimeAdded)
	if err != nil { return CounselRequest{}, err }
	return counselRequest, nil
}

func DeleteCounselRequest(db *sql.DB, userID int64) (CounselRequest, error) {
	_, err := db.Exec(`DELETE FROM wn_counsel_request WHERE user_id = $1`, userID)
	if err != nil { return CounselRequest{}, err }
	return CounselRequest{ UserID: userID }, nil
}

func AcceptCounselRequest(db *sql.DB, recipientID int64, providerID int64) (GroupWithUsers, error) {
	authorized := AuthoriseProvider(db, providerID)
	if !authorized { return GroupWithUsers{}, http_error.UnauthorizedError }
	present, err := CheckCounselRequest(db, recipientID)
	if err != nil { return GroupWithUsers{}, err }
	if !present { return GroupWithUsers{}, http_error.NotFoundError }
	user, err := GetUser(db, providerID)
	if err != nil { return GroupWithUsers{}, err }
	if !IsProvider(user) { return GroupWithUsers{}, http_error.UnauthorizedError }
	group := Group{
		GroupName: "Counsel Room",
		GroupDescription: "Welcome to your new Counsel Room",
		Category: "COUNSEL",
	}
	groupWithUsers, err := AddGroupWithUserIDs(db, group, []int64{providerID, recipientID})
	if err != nil { return GroupWithUsers{}, err }
	if _, fatal := DeleteCounselRequest(db, recipientID); fatal != nil {
		log.Fatal(fmt.Sprintf("Failed to remove counsel request after creating group. Fatal: %v", fatal))
	}
	return groupWithUsers, nil
}
