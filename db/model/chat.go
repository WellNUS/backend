package model

import (
	"time"
	"database/sql"
)

const (
	ServerUserID = -1

	MessageTag = 0
	ChatStatusTag = 1
)

// Message
type Message struct {
	UserID 		int64		`json:"user_id"`	
	GroupID		int64		`json:"group_id"`
	TimeAdded 	time.Time	`json:"time_added"`
	Msg			string		`json:"msg"`
}

type MessagePayload struct {
	Tag 		int 	`json:"tag"`
	SenderName	string	`json:"sender_name"`
	GroupName 	string	`json:"group_name"`
	Message		Message	`json:"message"`
}

type MessagesChunk struct {
	EarliestTime		time.Time			`json:"earliest_time"`
	LatestTime			time.Time 			`json:"latest_time"`
	MessagePayloads		[]MessagePayload 	`json:"message_payloads"`
}

func (m Message) IsServerMessage() bool {
	return m.UserID == ServerUserID;
}

func (m Message) Payload(db *sql.DB) (MessagePayload, error) {
	group, err := GetGroup(db, m.GroupID)
	if err != nil { return MessagePayload{}, err }
	var senderName string
	if m.IsServerMessage() {
		senderName = "[WellNUS Server]"
	} else {
		sender, err := GetUser(db, m.UserID)
		if err != nil { return  MessagePayload{}, err }
		senderName = sender.FirstName
	}
	return MessagePayload{Tag: MessageTag, SenderName: senderName, GroupName: group.GroupName, Message: m}, nil
}

// Chat Status
type ChatStatusPayload struct {
	Tag						int 			`json:"tag"`
	GroupID					int64			`json:"group_id"`
	GroupName				string			`json:"group_name"`
	SortedOnlineMembers 	[]User			`json:"sorted_online_users"`
	SortedInChatMembers		[]User			`json:"sorted_in_chat_users"`
}