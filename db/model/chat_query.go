package model

import (
	"database/sql"
	"time"
)

func readMessagePayloads(rows *sql.Rows) ([]MessagePayload, error) {
	messagePayloads := make([]MessagePayload, 0)
	for rows.Next() {
		var messagePayload MessagePayload
		err := rows.Scan(
			&messagePayload.SenderName,
			&messagePayload.GroupName,
			&messagePayload.Message.UserID,
			&messagePayload.Message.GroupID,
			&messagePayload.Message.TimeAdded,
			&messagePayload.Message.Msg)
		if err != nil { return nil, err }
		messagePayload.Tag = MessageTag
		messagePayloads = append(messagePayloads, messagePayload)
	}
	return messagePayloads, nil
}

func GetMessagesChunkOfGroupCustomise(db *sql.DB, groupID int64, latestTime time.Time, limit int64) (MessagesChunk, error) {
	var rows *sql.Rows
	var err error
	if limit <= 0 {
		rows, err = db.Query(
			`WITH t AS (
				SELECT 
					wn_user.first_name,
					wn_group.group_name,
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
					wn_user.first_name,
					wn_group.group_name,
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
	if err != nil { return MessagesChunk{}, err }
	messagePayloads, err := readMessagePayloads(rows)
	if err != nil { return MessagesChunk{}, err }

	messagesChunk := MessagesChunk{MessagePayloads: messagePayloads}
	if l := len(messagePayloads); l > 0 {
		messagesChunk.EarliestTime = messagePayloads[0].Message.TimeAdded
		messagesChunk.LatestTime = messagePayloads[len(messagePayloads) - 1].Message.TimeAdded
	}
	return messagesChunk, nil
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