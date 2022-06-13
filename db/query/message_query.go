package query

import (
	"wellnus/backend/model"
	"database/sql"
	"time"
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
			&loadedMessage.Group.ID,
			&loadedMessage.Group.GroupName,
			&loadedMessage.Group.GroupDescription,
			&loadedMessage.Group.Category,
			&loadedMessage.Group.OwnerID,
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
	group, err := getGroup(db, message.GroupID)
	if err != nil { return LoadedMessage{}, err }
	return LoadedMessage{ Sender: sender, Group: group, Message: message }, nil
}

func GetLoadedMessagesOfGroupCustomise(db *sql.DB, groupID int64, latestTime time.Time, limit int64) ([]LoadedMessage, error) {
	var rows *sql.Rows
	var err error
	if limit <= 0 {
		rows, err = db.Query(
			`WITH t AS (
				SELECT 
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
					wn_group.owner_id,
					wn_message.user_id,
					wn_message.group_id,
					wn_message.time_added,
					wn_message.msg
				FROM wn_message 
				JOIN wn_user
				ON wn_message.user_id = wn_user.id
				JOIN wn_group
				ON wn_message.group_id = wn_group.id
				WHERE wn_message.group_id = $1 AND wn_message.time_added < $2
				ORDER BY time_added DESC
			) SELECT * FROM t ORDER BY time_added ASC;`,
			groupID,
			latestTime)
	} else {
		rows, err = db.Query(
			`WITH t AS (
				SELECT 
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
					wn_group.owner_id,
					wn_message.user_id,
					wn_message.group_id,
					wn_message.time_added,
					wn_message.msg
				FROM wn_message 
				JOIN wn_user
				ON wn_message.user_id = wn_user.id
				JOIN wn_group
				ON wn_message.group_id = wn_group.id
				WHERE wn_message.group_id = $1 AND wn_message.time_added < $2
				ORDER BY time_added DESC
				LIMIT $3
			) SELECT * FROM t ORDER BY time_added ASC;`,
			groupID,
			latestTime,
			limit)
	}
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