package model

import (
	"time"
	"database/sql"
	"errors"
	"bytes"
)

const (
	MessageTag = 0
	ChatStatusTag = 1
	VideoTag = 2
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Chat Status
type ChatStatusPayload struct {
	Tag						int 			`json:"tag"`
	GroupID					int64			`json:"group_id"`
	GroupName				string			`json:"group_name"`
	SortedInChatMembers		[]User			`json:"sorted_in_chat_members"`
	SortedOnlineMembers 	[]User			`json:"sorted_online_members"`
	SortedOfflineMembers	[]User			`json:"sorted_offline_members"`
}

// SentData
type SentData struct {
	Tag 		int 		`json:"tag"`
	Client		*Client		`json:"client"`
	Data		interface{} `json:"data"`
}

func (sd SentData) ToMessage() (Message, error) {
	if (sd.Tag != MessageTag) {
		return Message{}, errors.New("sentData did not have message tag")
	}
	msg, ok := sd.Data.(string)
	if !ok {
		return Message{}, errors.New("data in sentData cannot be asserted as string")
	}
	msg = string(bytes.TrimSpace(bytes.Replace([]byte(msg), newline, space, -1)))
	return Message{ UserID: sd.Client.UserID, GroupID: sd.Client.GroupID, TimeAdded: time.Now(), Msg: msg }, nil
}

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

func (m Message) Payload(db *sql.DB) (MessagePayload, error) {
	group, err := GetGroup(db, m.GroupID)
	if err != nil { return MessagePayload{}, err }
	sender, err := GetUser(db, m.UserID)
	if err != nil { return  MessagePayload{}, err }
	senderName := sender.FirstName
	return MessagePayload{Tag: MessageTag, SenderName: senderName, GroupName: group.GroupName, Message: m}, nil
}