package model

import (
	"database/sql"
)

type ProviderSetting struct {
	UserID 			int64 		`json:"user_id"`
	Intro			string 		`json:"intro"`
	Specialities	[]string	`json:"specialities"`
}

type Provider struct {
	User 		User 			`json:"user"`
	Setting		ProviderSetting	`json:"setting"`
}

type ProviderWithEvents struct {
	Provider 	Provider 	`json:"provider"`
	Events		[]Event		`json:"events"`
}

func (ps ProviderSetting) LoadProviderSetting(db *sql.DB) (Provider, error) {
	user, err := GetUser(db, ps.UserID)
	if err != nil { return Provider{}, err }
	return Provider{ User: user, Setting: ps }, nil
}

func (p Provider) LoadProvider(db *sql.DB) (ProviderWithEvents, error) {
	events, err := GetAllEventsOfUser(db, p.User.ID)
	if err != nil { return ProviderWithEvents{}, err }
	return ProviderWithEvents{ Provider: p, Events: events }, nil
}