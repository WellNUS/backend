package model

import (
	"wellnus/backend/config"

	"database/sql"
	"math/rand"
	"time"
	"sort"
	"log"
	"bytes"
	"encoding/json"
	"os"
)

// Determines the compatibility between 2 loaded match request on a scale of 0 to 12
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

// Retrieves all loadedMatchRequest and performs matching on all available match request
func PerformMatching(db *sql.DB) ([]GroupWithUsers, error) {
	loadedMatchRequests, err := GetAllLoadedMatchRequest(db)
	if err != nil { return nil, err }
	if len(loadedMatchRequests) < config.MatchRequestThreshold { return make([]GroupWithUsers, 0), nil }

	groupsWithUsers := make([]GroupWithUsers, 0)
	group := Group{
		GroupName: "Support Group",
		GroupDescription: "Welcome to your new Support Group",
		Category: "SUPPORT",
	}

	for len(loadedMatchRequests) >= config.GroupSize {
		groupingIndices, remainingIndices := GetGroupingRemainingIndices(loadedMatchRequests)

		groupingUserIDs := make([]int64, len(groupingIndices))
		for i, index := range groupingIndices {
			userID := loadedMatchRequests[index].MatchRequest.UserID
			groupingUserIDs[i] = userID
			_, err := DeleteMatchRequestOfUser(db, userID)
			if err != nil { return nil, err }
		}

		groupWithUsers, err := AddGroupWithUserIDs(db, group, groupingUserIDs)
		if err != nil { return nil, err }
		groupsWithUsers = append(groupsWithUsers, groupWithUsers)

		remainingLMRs := make([]LoadedMatchRequest, len(remainingIndices))
		for i, index := range remainingIndices {
			remainingLMRs[i] = loadedMatchRequests[index]
		}
		loadedMatchRequests = remainingLMRs
	}
	// writeToStdOut(groupsWithUsers)
	return groupsWithUsers, nil
}

// Return a slice of the indices of match request that will form a group among given loadmatchrequests
// Done by randomly selecting a pivoting match request and building the group to best suit that pivot
func GetGroupingRemainingIndices(lmrs []LoadedMatchRequest) ([]int, []int) {
	l := len(lmrs)
	if l < config.GroupSize { return nil, nil }
	p := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(l)
	scoreMap := make([]int, l)
	indices := make([]int, l)
	for j, lmr := range lmrs {
		scoreMap[j] = Compatibility(lmrs[p], lmr)
		indices[j] = j
	}
	sort.Slice(indices, func(i, j int) bool {
		return scoreMap[indices[i]] > scoreMap[indices[j]]
	})
	return indices[:config.GroupSize], indices[config.GroupSize:]
}

func writeToStdOut(item interface{}) {
	b, err := json.Marshal(item)
	if err != nil {
		log.Fatal(err)
	}
	var out bytes.Buffer
	json.Indent(&out, b, "=", " ")
	out.WriteTo(os.Stdout)
}