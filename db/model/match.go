package model

import (
	"wellnus/backend/config"

	"time"
	"database/sql"
	"errors"
	"sort"
	"fmt"
)

// Helper for matching requests
func Compatibility(loadedMatchRequest1, loadedMatchRequest2 LoadedMatchRequest) int {
	totalScore := 0

	// Faculty check max 4 points
	fac1, fac2 := loadedMatchRequest1.User.Faculty, loadedMatchRequest2.User.Faculty
	pfac1, pfac2 := loadedMatchRequest1.MatchSetting.FacultyPreference, loadedMatchRequest2.MatchSetting.FacultyPreference
	if fac1 == fac2 {
		if pfac1 == "SAME" && pfac2 == "SAME" {
			totalScore += 4
		} else if pfac1 == "MIX" && pfac2 == "MIX" {
			totalScore += 0
		} else {
			totalScore += 2
		}
	} else {
		if pfac1 == "SAME" && pfac2 == "SAME" {
			totalScore += 0
		} else if pfac1 == "MIX" && pfac2 == "MIX" {
			totalScore += 4
		} else {
			totalScore += 2
		}
	}

	// MBTI check max 4 points
	if (loadedMatchRequest1.MatchSetting.MBTI[0] == loadedMatchRequest2.MatchSetting.MBTI[0]) { totalScore += 1 }
	if (loadedMatchRequest1.MatchSetting.MBTI[1] == loadedMatchRequest2.MatchSetting.MBTI[1]) { totalScore += 1 }
	if (loadedMatchRequest1.MatchSetting.MBTI[2] == loadedMatchRequest2.MatchSetting.MBTI[2]) { totalScore += 1 }
	if (loadedMatchRequest1.MatchSetting.MBTI[3] == loadedMatchRequest2.MatchSetting.MBTI[3]) { totalScore += 1 }

	// Hobbies check max 4 points (max of 4 hobbies)
	for _, hobby1 := range loadedMatchRequest1.MatchSetting.Hobbies {
		for _, hobby2 := range loadedMatchRequest2.MatchSetting.Hobbies {
			if hobby1 == hobby2 {
				totalScore += 1
				break
			}
		}
	}
	return totalScore // out of 12
}

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

type SortedPairs [][2]int

type LoadedMatchRequests []LoadedMatchRequest

func (mr MatchRequest) LoadMatchRequest(db *sql.DB) (LoadedMatchRequest, error) {
	user, err := GetUser(db, mr.UserID)
	if err != nil { return LoadedMatchRequest{}, err }
	matchSetting, err := GetMatchSettingOfUser(db, mr.UserID)
	if err != nil { return LoadedMatchRequest{}, err }
	return LoadedMatchRequest{ MatchRequest: mr, User: user, MatchSetting: matchSetting }, nil
}

func (sortedPairs SortedPairs) MakeMore() func(int, int) bool {
	return func(i, j int) bool {
		return sortedPairs[i][0] > sortedPairs[j][0]
	}
}


func (lmr LoadedMatchRequests) GetMostCompatible() (LoadedMatchRequests, LoadedMatchRequests, error) {
	fmt.Printf("%v \n", lmr)
	l := len(lmr)
	if l < config.GroupSizes { return nil, nil, errors.New("Insufficient requests to get most compatible group") }

	ph := make(map[int]SortedPairs)
	for i := 0; i < l; i++ {
		ph[i] = make(SortedPairs, 0)
	}
	for i := 0; i < l - 1; i++ {
		for j := i; j < l; j++ {
			c := Compatibility(lmr[i], lmr[j])
			ph[i] = append(ph[i], [2]int{c, j})
			ph[j] = append(ph[j], [2]int{c, i})
		} 
	}
	for i := 0; i < l; i++ {
		sort.Slice(ph[i], ph[i].MakeMore())
		fmt.Printf("%d: %v \n", i, ph[i])
	}

	maxScore := -1
	bestPivot := -1
	for i:=0; i < l; i++ {
		score := 0
		for _, pair := range ph[i][:config.GroupSizes-1] {
			score += pair[0]
		}
		if score > maxScore {
			maxScore = score
			bestPivot = i
		}
	}

	topNew := make(LoadedMatchRequests, 0)
	topRemaining := make(LoadedMatchRequests, 0)
	for _, pair := range ph[bestPivot][:config.GroupSizes] {
		topNew = append(topNew, lmr[pair[1]])
	}
	for _, pair := range ph[bestPivot][config.GroupSizes:] {
		topRemaining = append(topRemaining, lmr[pair[1]])
	}
	return topNew, topRemaining, nil
}