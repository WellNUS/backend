package model

import (
	"time"
	"database/sql"
)

// Direct
type DirectMessage struct {
	SenderID 	int64		`json:"sender_id"`
	RecipientID int64		`json:"recipient_id"`
	TimeAdded 	time.Time	`json:"time_added"`
	Msg			string		`json:"msg"`
}

type DirectMessagePayload struct {
	Tag 			int				`json:"tag"`
	Label			string			`json:"label"`
	SenderName 		string			`json:"sender_name"`
	RecipientName	string			`json:"recipient_name"`
	Message			DirectMessage	`json:"message"`
}

type DirectMessagesChunk struct {
	EarliestTime	time.Time					`json:"earliest_time"`
	LatestTime		time.Time					`json:"latest_time"`
	Payloads 		[]DirectMessagePayload		`json:"payloads"`
}

func (d DirectMessage) Payload(db *sql.DB) (DirectMessagePayload, error) {
	sender, err := GetUser(db, d.SenderID)
	if err != nil { return DirectMessagePayload{}, err }
	recipient, err := GetUser(db, d.RecipientID)
	if err != nil { return DirectMessagePayload{}, err }
	return DirectMessagePayload{
		Tag: MessageTag,
		Label: DirectMessageLabel,
		SenderName: sender.FirstName,
		RecipientName: recipient.FirstName,
		Message: d,
	}, nil
}


// Message
type GroupMessage struct {
	UserID 		int64		`json:"user_id"`	
	GroupID		int64		`json:"group_id"`
	TimeAdded 	time.Time	`json:"time_added"`
	Msg			string		`json:"msg"`
}

type GroupMessagePayload struct {
	Tag 		int 			`json:"tag"`
	Label		string			`json:"label"`
	SenderName	string			`json:"sender_name"`
	GroupName 	string			`json:"group_name"`
	Message		GroupMessage	`json:"message"`
}

type GroupMessagesChunk struct {
	EarliestTime		time.Time				`json:"earliest_time"`
	LatestTime			time.Time 				`json:"latest_time"`
	Payloads			[]GroupMessagePayload 	`json:"payloads"`
}

func (m GroupMessage) Payload(db *sql.DB) (GroupMessagePayload, error) {
	group, err := GetGroup(db, m.GroupID)
	if err != nil { return GroupMessagePayload{}, err }
	sender, err := GetUser(db, m.UserID)
	if err != nil { return  GroupMessagePayload{}, err }
	return GroupMessagePayload{
		Tag: MessageTag,
		Label: GroupMessageLabel,
		SenderName: sender.FirstName, 
		GroupName: group.GroupName, 
		Message: m,
	}, nil
}