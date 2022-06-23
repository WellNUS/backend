package model

import (
	"time"
	"database/sql"
)

type MatchSetting struct {
	UserID 				int64 		`json:"user_id"`
	FacultyPreference 	string 		`json:"faculty_preference"`
	Hobbies 			[]string	`json:"hobbies"`
	MBTI				string		`json:"mbti"`
}

type MatchRequest struct {
	UserID 		int64		`json:"user_id"`
	TimeAdded	time.Time	`json:"time_added"`
}

type LoadedMatchRequest struct {
	MatchRequest 	MatchRequest 	`json:"match_request"`
	User			User			`json:"user"`
	MatchSetting	MatchSetting	`json:"match_setting"`
}

func (mr MatchRequest) LoadMatchRequest(db *sql.DB) (LoadedMatchRequest, error) {
	user, err := GetUser(db, mr.UserID)
	if err != nil { return LoadedMatchRequest{}, err }
	matchSetting, err := GetMatchSettingOfUser(db, mr.UserID)
	if err != nil { return LoadedMatchRequest{}, err }
	return LoadedMatchRequest{ MatchRequest: mr, User: user, MatchSetting: matchSetting }, nil
}