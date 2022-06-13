package model

import (
	"time"
)

type User struct {
	ID 				int64 	`json:"id"`
	FirstName 		string 	`json:"first_name"`
	LastName 		string	`json:"last_name"`
	Gender			string 	`json:"gender"`
	Faculty			string 	`json:"faculty"`
	Email			string	`json:"email"`
	UserRole		string 	`json:"user_role"`
	Password		string 	`json:"password"`
	PasswordHash 	string	`json:"password_hash"`
}

type UserWithGroups struct {
	User 	User 	`json:"user"`
	Groups 	[]Group `json:"groups"`
}

type Group struct {
	ID					int64	`json:"id"`
	GroupName			string	`json:"group_name"`
	GroupDescription 	string	`json:"group_description"`
	Category			string 	`json:"category"`
	OwnerID				int64	`json:"owner_id"`
}

type GroupWithUsers struct {
	Group	Group	`json:"group"`
	Users	[]User	`json:"users"`
}

type JoinRequestRespond struct {
	Approve bool `json:"approve"`
}

type JoinRequest struct {
	ID 				int64 	`json:"id"`
	UserID 			int64 	`json:"user_id"`
	GroupID 		int64 	`json:"group_id"`
}

type LoadedJoinRequest struct {
	JoinRequest		JoinRequest 	`json:"join_request"`
	User			User			`json:"user"`
	Group			Group			`json:"group"`
}

type Message struct {
	UserID 		int64		`json:"user_id"`
	GroupID		int64		`json:"group_id"`
	TimeAdded 	time.Time	`json:"time_added"`
	Msg			string		`json:"msg"`
}

type LoadedMessage struct {
	Sender 			User	`json:"sender"`
	Group			Group	`json:"group"`
	Message			Message	`json:"message"`
}

type LoadedMessagesPacket struct {
	EarliestTime		time.Time			`json:"earliest_time"`
	LatestTime			time.Time 			`json:"latest_time"`
	LoadedMessages 		[]LoadedMessage 	`json:"loaded_messages"`
}