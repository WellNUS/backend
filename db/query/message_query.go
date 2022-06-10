package query

import (
	"wellnus/backend/db/model"
	"database/sql"
)

type Message = model.Message
type LoadedMessage = model.LoadedMessage

func LoadMessage(db *sql.DB, message Message) (LoadedMessage, error) {
	owner, err := getUser(db, message.UserID)
	if err != nil { return LoadedMessage{}, err }
	users, err := GetAllUsersOfGroup(db, message.GroupID)
	if err != nil { return LoadedMessage{}, err }
	return LoadedMessage{ Owner: owner, Users: users, Message: message }, nil
}

/*
AddMessage(db *sql.DB, message Message) (Message, error) {

}
*/