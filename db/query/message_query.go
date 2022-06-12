package query

import (
	"wellnus/backend/db/model"
	"database/sql"

	"fmt"
)

type Message = model.Message
type LoadedMessage = model.LoadedMessage

func readLoadedMessages(rows *sql.Rows) ([]LoadedMessage, error) {
	loadedMessages := make([]LoadedMessage, 0)
	for rows.Next() {
		var loadedMessage LoadedMessage
		err := rows.Scan(
			&loadedMessage.Sender.ID,
			&loadedMessage.Sender.FirstName,
			&loadedMessage.Sender.LastName,
			&loadedMessage.Sender.Gender,
			&loadedMessage.Sender.Faculty,
			&loadedMessage.Sender.Email,
			&loadedMessage.Sender.UserRole,
			&loadedMessage.Sender.PasswordHash,
			&loadedMessage.Message.UserID,
			&loadedMessage.Message.GroupID,
			&loadedMessage.Message.TimeAdded,
			&loadedMessage.Message.Msg)
		if err != nil { return nil, err }
		loadedMessages = append(loadedMessages, loadedMessage)
	}
	return loadedMessages, nil
}

func LoadMessage(db *sql.DB, message Message) (LoadedMessage, error) {
	sender, err := getUser(db, message.UserID)
	if err != nil { return LoadedMessage{}, err }
	return LoadedMessage{ Sender: sender, Message: message }, nil
}

func GetAllLoadedMessagesOfGroup(db *sql.DB, groupID int64) ([]LoadedMessage, error) {
	query := fmt.Sprintf(
		`SELECT 
			wn_user.id,
			wn_user.first_name,
			wn_user.last_name,
			wn_user.gender,
			wn_user.faculty,
			wn_user.email,
			wn_user.user_role,
			wn_user.password_hash,
			wn_message.user_id,
			wn_message.group_id,
			wn_message.time_added,
			wn_message.msg
		FROM wn_message 
		JOIN wn_user
		ON wn_message.user_id = wn_user.id
		WHERE wn_message.group_id = %d
		ORDER BY time_added ASC;`,
		groupID)
	rows, err := db.Query(query)
	if err != nil { return nil, err }
	loadedMessages, err := readLoadedMessages(rows)
	if err != nil { return nil, err }
	return loadedMessages, nil
}

func AddMessage(db *sql.DB, message Message) error {
	_, err := db.Exec(
		`INSERT INTO wn_message (
			user_id,
			group_id,
			time_added,
			msg
		) values ($1, $2, $3, $4)`,
		message.UserID,
		message.GroupID,
		message.TimeAdded,
		message.Msg)
	return err
}