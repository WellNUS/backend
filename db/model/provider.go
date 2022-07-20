package model

import (
	"database/sql"
)

type ProviderSetting struct {
	UserID 			int64 		`json:"user_id"`
	Intro			string 		`json:"intro"`
	Specialities	[]string	`json:"specialities"`
}

type ProviderWithSetting struct {
	User 		User 			`json:"user"`
	Setting		ProviderSetting	`json:"setting"`
}

func (ps ProviderSetting) LoadProviderSettings(db *sql.DB) (ProviderWithSetting, error) {
	user, err := GetUser(db, ps.UserID)
	if err != nil { return ProviderWithSetting{}, err }
	return ProviderWithSetting{ User: user, Setting: ps }, nil
}