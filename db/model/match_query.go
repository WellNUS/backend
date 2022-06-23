package model

import (
	"wellnus/backend/router/http_helper/http_error"
	"database/sql"
	"time"
	"github.com/lib/pq"
)

func ReadMatchSettings(rows *sql.Rows) ([]MatchSetting, error) {
	matchSettings := make([]MatchSetting, 0)
	for rows.Next() {
		var matchSetting MatchSetting
		if err := rows.Scan(
			&matchSetting.UserID,
			&matchSetting.FacultyPreference,
			pq.Array(&matchSetting.Hobbies),
			&matchSetting.MBTI);
			err != nil {
				return nil, err
			}
		matchSettings = append(matchSettings, matchSetting)
	}
	return matchSettings, nil
}

func ReadMatchRequests(rows *sql.Rows) ([]MatchRequest, error) {
	matchRequests := make([]MatchRequest, 0)
	for rows.Next() {
		var matchRequest MatchRequest
		if err := rows.Scan(&matchRequest.UserID, &matchRequest.TimeAdded); err != nil {
			return nil, err
		}
		matchRequests = append(matchRequests, matchRequest)
	}
	return matchRequests, nil
}

func ReadLoadedMatchRequests(rows *sql.Rows) ([]LoadedMatchRequest, error) {
	loadedMatchRequests := make([]LoadedMatchRequest, 0)
	for rows.Next() {
		var loadedMatchRequest LoadedMatchRequest
		if err := rows.Scan(
			&loadedMatchRequest.MatchRequest.UserID,
			&loadedMatchRequest.MatchRequest.TimeAdded,
			&loadedMatchRequest.User.ID,
			&loadedMatchRequest.User.FirstName,
			&loadedMatchRequest.User.LastName,
			&loadedMatchRequest.User.Gender,
			&loadedMatchRequest.User.Faculty,
			&loadedMatchRequest.User.Email,
			&loadedMatchRequest.User.UserRole,
			&loadedMatchRequest.User.PasswordHash,
			&loadedMatchRequest.MatchSetting.UserID,
			&loadedMatchRequest.MatchSetting.FacultyPreference,
			pq.Array(&loadedMatchRequest.MatchSetting.Hobbies),
			&loadedMatchRequest.MatchSetting.MBTI);
			err != nil {
				return nil, err
			}
		loadedMatchRequests = append(loadedMatchRequests, loadedMatchRequest)
	}
	return loadedMatchRequests, nil
}

// Match setting

func GetMatchSettingOfUser(db *sql.DB, userID int64) (MatchSetting, error){
	rows, err := db.Query(`SELECT * FROM wn_match_setting WHERE user_id = $1`, userID)
	if err != nil { return MatchSetting{}, err }
	matchSettings, err := ReadMatchSettings(rows);
	if err != nil { return MatchSetting{}, err }
	if len(matchSettings) == 0 { return MatchSetting{}, http_error.NotFoundError }
	return matchSettings[0], nil
}

func AddUpdateMatchSettingOfUser(db *sql.DB, matchSetting MatchSetting, userID int64) (MatchSetting, error) {
	matchSetting.UserID = userID
	_, err := db.Exec(
		`INSERT INTO wn_match_setting (
			user_id,
			faculty_preference,
			hobbies,
			mbti
		) VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			faculty_preference = EXCLUDED.faculty_preference,
			hobbies = EXCLUDED.hobbies,
			mbti = EXCLUDED.mbti`,
		matchSetting.UserID,
		matchSetting.FacultyPreference,
		pq.Array(matchSetting.Hobbies),
		matchSetting.MBTI)
	if err != nil { return MatchSetting{}, err }
	return matchSetting, nil
}

func DeleteMatchSettingOfUser(db *sql.DB, userID int64) (MatchSetting, error) {
	if _, err := db.Exec(`DELETE FROM wn_match_setting WHERE user_id = $1`, userID); err != nil {
		return MatchSetting{}, err
	}
	return MatchSetting{ UserID: userID }, nil
}

// Match Request

func GetMatchRequestCount(db *sql.DB) (int64, error) {
	rows, err := db.Query(`SELECT COUNT(*) FROM wn_match_request`)
	if err != nil { return 0, err }
	rows.Next()
	var count int64
	if err := rows.Scan(&count); err != nil { return 0, err }
	return count, nil
}

func GetMatchRequestOfUser(db *sql.DB, userID int64) (MatchRequest, error) {
	rows, err := db.Query(`SELECT * FROM wn_match_request WHERE user_id = $1`, userID)
	if err != nil { return MatchRequest{}, err }
	matchRequests, err := ReadMatchRequests(rows)
	if err != nil { return MatchRequest{}, err }
	if len(matchRequests) == 0 { return MatchRequest{}, http_error.NotFoundError }
	return matchRequests[0], nil
}

func GetLoadedMatchRequestOfUser(db *sql.DB, userID int64) (LoadedMatchRequest, error) {
	matchRequest, err := GetMatchRequestOfUser(db, userID)
	if err != nil { return LoadedMatchRequest{}, err }
	loadedMatchRequest, err := matchRequest.LoadMatchRequest(db)
	if err != nil { return LoadedMatchRequest{}, err }
	return loadedMatchRequest, nil
}

func GetAllLoadedMatchRequest(db *sql.DB) ([]LoadedMatchRequest, error) {
	rows, err := db.Query(
		`SELECT
			wn_match_request.user_id,
			wn_match_request.time_added,
			wn_user.id,
			wn_user.first_name,
			wn_user.last_name,
			wn_user.gender,
			wn_user.faculty,
			wn_user.email,
			wn_user.user_role,
			wn_user.password_hash,
			wn_match_setting.user_id,
			wn_match_setting.faculty_preference,
			wn_match_setting.hobbies,
			wn_match_setting.mbti
		FROM wn_match_request
		JOIN wn_user ON wn_match_request.user_id = wn_user.id
		LEFT JOIN wn_match_setting ON wn_match_request.user_id = wn_match_setting.user_id`)
	if err != nil { return nil, err }
	loadedMatchRequests, err := ReadLoadedMatchRequests(rows)
	if err != nil { return nil, err }
	return loadedMatchRequests, nil
}

// Continue with add match request and match algorithm after a threshold is met
func AddMatchRequest(db *sql.DB, userID int64) (MatchRequest, error) {
	matchRequest := MatchRequest{ UserID: userID, TimeAdded: time.Now() }
	_, err := db.Exec(
		`INSERT INTO wn_match_request (
			user_id,
			time_added
		) VALUES ($1, $2)`,
		matchRequest.UserID,
		matchRequest.TimeAdded)
	if err != nil { return MatchRequest{}, err }
	PerformMatching(db)
	return matchRequest, nil
}

func DeleteMatchRequestOfUser(db *sql.DB, userID int64) (MatchRequest, error) {
	if _, err := db.Exec(`DELETE FROM wn_match_request WHERE user_id = $1`, userID); err != nil {
		return MatchRequest{}, err
	}
	return MatchRequest{ UserID: userID }, nil
} 