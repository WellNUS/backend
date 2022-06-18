package model

import (
	"wellnus/backend/router/misc/http_error"
	"database/sql"
	"math/rand"
	"time"
)

// Helper function
func readSessions(rows *sql.Rows) ([]Session, error) {
	sessions := make([]Session, 0)
	for rows.Next() {
		var session Session
		if err := rows.Scan(&session.SessionKey, &session.UserID); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func GenerateNewSessionKey() string {
	Rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, SessionKeyLength)
	charSetLen := len(SessionKeyCharSet)
	for i := range b {
		b[i] = SessionKeyCharSet[Rand.Intn(charSetLen)]
	}
	return string(b)
}


// Main function
func GetUserIDFromSessionKey(db *sql.DB, sessionKey string) (int64, error) {
	rows, err := db.Query(`SELECT * FROM wn_session WHERE session_key = $1`, sessionKey)
	if err != nil { return 0, err }
	sessions, err := readSessions(rows)
	if err != nil { return 0, err }
	if len(sessions) == 0 { return 0, http_error.UnauthorizedError }
	return sessions[0].UserID, nil
}

func DeleteSessionWithSessionKey(db *sql.DB, sessionKey string) error {
	_, err := db.Exec(`DELETE FROM wn_session WHERE session_key = $1`, sessionKey)
	return err
}

func CreateNewSession(db *sql.DB, userID int64) (string, error) {
	newSessionKey := GenerateNewSessionKey()
	_, err := db.Exec(
		`INSERT INTO wn_session (
			session_key,
			user_id
		) VALUES ($1, $2)
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			session_key = EXCLUDED.session_key,
			user_id = EXCLUDED.user_id`,
		newSessionKey,
		userID)
	if err != nil { return "", err }
	return newSessionKey, nil
}

