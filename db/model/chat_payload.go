package model

import (
	"time"
	"errors"
	"bytes"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	MessageTag = 0
	ChatStatusTag = 1 //Cant be sent by client
	VideoTag = 2
)

const (
	GroupMessageLabel = "GROUP"
	DirectMessageLabel = "DIRECT"
	ServerMessageLabel = "SERVER" //Cant be sent by client
)

const (
	InChatStatus = "In Chat"
	OnlineStatus = "Online"
	OfflineStatus = "Offline"
)

const (
	GroupStatusLabel = "GROUP"
	UserStatusLabel = "USER"
)

// SentData
// SentData could either be message, chatstatus or video
// When sending to the server, sentData must be tag appropriately
// Similarly when sending back to client, the data will be tagged appropriately
type SentData struct {
	Tag 		int 		`json:"tag"`
	Client		*Client		`json:"client"`
	Data		interface{} `json:"data"`
}

type ServerMessagePayload struct {
	Tag			int		`json:"tag"`
	Label		string	`json:"label"`
	Msg			string	`json:"msg"`
}

func (sd SentData) ToGroupMessage() (GroupMessage, error) {
	if sd.Tag != MessageTag {
		return GroupMessage{}, errors.New("sentData did not have message tag")
	}
	if !sd.Client.TargetIsGroup {
		return GroupMessage{}, errors.New("sentData is not a group message")
	}
	msg, ok := sd.Data.(string)
	if !ok {
		return GroupMessage{}, errors.New("data in sentData cannot be asserted as string")
	}
	msg = string(bytes.TrimSpace(bytes.Replace([]byte(msg), newline, space, -1)))
	return GroupMessage{
		UserID: sd.Client.UserID, 
		GroupID: sd.Client.TargetID, 
		TimeAdded: time.Now(), 
		Msg: msg,
	}, nil
}

func (sd SentData) ToDirectMessage() (DirectMessage, error) {
	if sd.Tag != MessageTag {
		return DirectMessage{}, errors.New("sentData did not have message tag")
	}
	if sd.Client.TargetIsGroup {
		return DirectMessage{}, errors.New("sentData is not a direct message")
	}
	msg, ok := sd.Data.(string)
	if !ok {
		return DirectMessage{}, errors.New("data in sentData cannot be asserted as string")
	}
	msg = string(bytes.TrimSpace(bytes.Replace([]byte(msg), newline, space, -1)))
	return DirectMessage{ 
		SenderID: sd.Client.UserID, 
		RecipientID: sd.Client.TargetID, 
		TimeAdded: time.Now(), 
		Msg: msg,
	}, nil
}

// Chat Status
// Members are in only 1 of 3 states (in chat, online or offline)
// inChat means the member is on the given chat page
// online means the member is connected but on some other chat page
// offline means the member is not connected
type GroupStatusPayload struct {
	Tag						int 			`json:"tag"`
	Label					string			`json:"label"`
	GroupID					int64			`json:"group_id"`
	GroupName				string			`json:"group_name"`
	SortedInChatMembers		[]User			`json:"sorted_in_chat_members"`
	SortedOnlineMembers 	[]User			`json:"sorted_online_members"`
	SortedOfflineMembers	[]User			`json:"sorted_offline_members"`
}

type UserStatusPayload struct {
	Tag		int 		`json:"tag"`
	Label 	string		`json:"label"`
	User	User		`json:"user"`
	Status	string		`json:"status"`
}