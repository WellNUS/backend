package model

const (
	SessionKeyLength = 128
	SessionKeyCharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

type Session struct {
	SessionKey 	string 	`json:"session_key"`
	UserID		int64	`json:"user_id"`	
}

type SessionResponse struct {
	LoggedIn 	bool 	`json:"logged_in"`
	User	 	User 	`json:"user"`
	SessionKey 	string 	`json:"session_key"`
}