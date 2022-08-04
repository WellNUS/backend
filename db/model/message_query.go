package model

import (
	"database/sql"
	"time"
)

func readDirectMessagePayloads(rows *sql.Rows) ([]DirectMessagePayload, error) {
	directMessagePayloads := make([]DirectMessagePayload, 0)
	for rows.Next() {
		var directMessagePayload DirectMessagePayload
		err := rows.Scan(
			&directMessagePayload.SenderName,
			&directMessagePayload.RecipientName,
			&directMessagePayload.Message.SenderID,
			&directMessagePayload.Message.RecipientID,
			&directMessagePayload.Message.TimeAdded,
			&directMessagePayload.Message.Msg)
		if err != nil { return nil, err }
		directMessagePayload.Tag = MessageTag
		directMessagePayload.Label = DirectMessageLabel
		directMessagePayloads = append(directMessagePayloads, directMessagePayload)
	}
	return directMessagePayloads, nil
}

func readGroupMessagePayloads(rows *sql.Rows) ([]GroupMessagePayload, error) {
	groupMessagePayloads := make([]GroupMessagePayload, 0)
	for rows.Next() {
		var groupMessagePayload GroupMessagePayload
		err := rows.Scan(
			&groupMessagePayload.SenderName,
			&groupMessagePayload.GroupName,
			&groupMessagePayload.Message.UserID,
			&groupMessagePayload.Message.GroupID,
			&groupMessagePayload.Message.TimeAdded,
			&groupMessagePayload.Message.Msg)
		if err != nil { return nil, err }
		groupMessagePayload.Tag = MessageTag
		groupMessagePayload.Label = GroupMessageLabel
		groupMessagePayloads = append(groupMessagePayloads, groupMessagePayload)
	}
	return groupMessagePayloads, nil
}

func GetDirectMessagesChunk(db *sql.DB, userID0 int64, userID1 int64, latestTime time.Time, limit int64) (DirectMessagesChunk, error) {
	var rows *sql.Rows
	var err error
	if limit <= 0 {
		rows, err = db.Query(
			`WITH t AS (
				SELECT 
					sender.first_name,
					recipient.first_name,
					wn_direct.sender_id,
					wn_direct.recipient_id,
					wn_direct.time_added,
					wn_direct.msg
				FROM wn_direct 
				JOIN wn_user AS sender
				ON wn_direct.sender_id = sender.id
				JOIN wn_user AS recipient 
				ON wn_direct.recipient_id = recipient.id
				WHERE (wn_direct.recipient_id = $1 AND wn_direct.sender_id = $2) 
				OR (wn_direct.sender_id = $1 AND wn_direct.recipient_id = $2) 
				WHERE wn_direct.time_added < $3
				ORDER BY time_added DESC
			) SELECT * FROM t ORDER BY time_added ASC;`,
			userID0,
			userID1,
			latestTime)
	} else {
		rows, err = db.Query(
			`WITH t AS (
				SELECT 
					sender.first_name,
					recipient.first_name,
					wn_direct.sender_id,
					wn_direct.recipient_id,
					wn_direct.time_added,
					wn_direct.msg
				FROM wn_direct 
				JOIN wn_user AS sender
				ON wn_direct.sender_id = sender.id
				JOIN wn_user AS recipient 
				ON wn_direct.recipient_id = recipient.id
				WHERE ((wn_direct.recipient_id = $1 AND wn_direct.sender_id = $2) 
				OR (wn_direct.sender_id = $1 AND wn_direct.recipient_id = $2))
				AND wn_direct.time_added < $3
				ORDER BY time_added DESC
				LIMIT $4
			) SELECT * FROM t ORDER BY time_added ASC;`,
			userID0,
			userID1,
			latestTime,
			limit)
	}
	if err != nil { return DirectMessagesChunk{}, err }
	defer rows.Close()
	directMessagesPayloads, err := readDirectMessagePayloads(rows)
	if err != nil { return DirectMessagesChunk{}, err }

	directMessagesChunk := DirectMessagesChunk{Payloads: directMessagesPayloads}
	if l := len(directMessagesPayloads); l > 0 {
		directMessagesChunk.EarliestTime = directMessagesPayloads[0].Message.TimeAdded
		directMessagesChunk.LatestTime = directMessagesPayloads[len(directMessagesPayloads) - 1].Message.TimeAdded
	}
	return directMessagesChunk, nil
}

func AddDirectMessage(db *sql.DB, directMessage DirectMessage) error {
	_, err := db.Exec(
		`INSERT INTO wn_direct (
			sender_id,
			recipient_id,
			time_added,
			msg
		) values ($1, $2, $3, $4)`,
		directMessage.SenderID,
		directMessage.RecipientID,
		directMessage.TimeAdded,
		directMessage.Msg)
	return err
}

func GetGroupMessagesChunk(db *sql.DB, groupID int64, latestTime time.Time, limit int64) (GroupMessagesChunk, error) {
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
	if err != nil { return GroupMessagesChunk{}, err }
	defer rows.Close()
	groupMessagePayloads, err := readGroupMessagePayloads(rows)
	if err != nil { return GroupMessagesChunk{}, err }

	groupMessagesChunk := GroupMessagesChunk{Payloads: groupMessagePayloads}
	if l := len(groupMessagePayloads); l > 0 {
		groupMessagesChunk.EarliestTime = groupMessagePayloads[0].Message.TimeAdded
		groupMessagesChunk.LatestTime = groupMessagePayloads[len(groupMessagePayloads) - 1].Message.TimeAdded
	}
	return groupMessagesChunk, nil
}

func AddGroupMessage(db *sql.DB, groupMessage GroupMessage) error {
	_, err := db.Exec(
		`INSERT INTO wn_message (
			user_id,
			group_id,
			time_added,
			msg
		) values ($1, $2, $3, $4)`,
		groupMessage.UserID,
		groupMessage.GroupID,
		groupMessage.TimeAdded,
		groupMessage.Msg)
	return err
}