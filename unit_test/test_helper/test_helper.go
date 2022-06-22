package test_helper

import (
	"wellnus/backend/db/model"

	"database/sql"
	"net/http"
	"net/http/httptest"
	"bytes"
	"errors"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

type User = model.User
type UserWithGroups = model.UserWithGroups
type SessionResponse = model.SessionResponse
type Group = model.Group
type GroupWithUsers = model.GroupWithUsers
type JoinRequest = model.JoinRequest
type LoadedJoinRequest = model.LoadedJoinRequest
type JoinRequestRespond = model.JoinRequestRespond
type MatchSetting = model.MatchSetting
type MatchRequest = model.MatchRequest
type LoadedMatchRequest = model.LoadedMatchRequest


func GetBufferFromRecorder(w *httptest.ResponseRecorder) *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Result().Body)
	return buf
}

func GetCookieFromRecorder(w *httptest.ResponseRecorder, name string) string {
	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}

func SimulateRequest(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func ReadInt(row *sql.Rows) (int, error) {
	row.Next()
	var c int
	if err := row.Scan(&c); err != nil { return 0, err }
	return c, nil
}

func GetUserFromRecorder(w *httptest.ResponseRecorder) (User, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return User{}, errors.New(buf.String())
	}

	// fmt.Printf("Response Body: %v \n", buf)
	var user User
	err := json.NewDecoder(buf).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUserWithGroupsFromRecorder(w *httptest.ResponseRecorder) (UserWithGroups, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return UserWithGroups{}, errors.New(buf.String())
	}

	// fmt.Printf("Response Body: %v \n", buf)
	var userWithGroups UserWithGroups
	err := json.NewDecoder(buf).Decode(&userWithGroups)
	if err != nil {
		return UserWithGroups{}, err
	}
	return userWithGroups, nil
}

func GetUsersFromRecorder(w *httptest.ResponseRecorder) ([]User, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return nil, errors.New(buf.String())
	}
	// fmt.Printf("Response Body: %v \n", buf)
	users := make([]User, 0)
	err := json.NewDecoder(buf).Decode(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetSessionResponseFromRecorder(w *httptest.ResponseRecorder) (SessionResponse, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Result().Body)
	if w.Code != http.StatusOK {
		return SessionResponse{}, errors.New(buf.String())
	}

	//fmt.Printf("Response Body: %v \n", buf)
	var sessionResponse SessionResponse
	err := json.NewDecoder(buf).Decode(&sessionResponse)
	if err != nil {
		return SessionResponse{}, err
	}
	return sessionResponse, nil
}

func GetGroupFromRecorder(w *httptest.ResponseRecorder) (Group, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return Group{}, errors.New(buf.String())
	}

	// fmt.Printf("Response Body: %v \n", buf)
	var group Group
	err := json.NewDecoder(buf).Decode(&group)
	if err != nil {
		return Group{}, err
	}
	return group, nil
}

func GetGroupsFromRecorder(w *httptest.ResponseRecorder) ([]Group, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return nil, errors.New(buf.String())
	}

	// fmt.Printf("Response Body: %v \n", buf)
	var groups []Group
	err := json.NewDecoder(buf).Decode(&groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func GetGroupWithUsersFromRecorder(w *httptest.ResponseRecorder) (GroupWithUsers, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return GroupWithUsers{}, errors.New(buf.String())
	}
	// fmt.Printf("Response Body: %v \n", buf)
	var groupWithUsers GroupWithUsers
	err := json.NewDecoder(buf).Decode(&groupWithUsers)
	if err != nil {
		return GroupWithUsers{}, err
	}
	return groupWithUsers, nil
}

func GetGroupsWithUsersFromRecorder(w *httptest.ResponseRecorder) ([]GroupWithUsers, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return nil, errors.New(buf.String())
	}
	// fmt.Printf("Response Body: %v \n", buf)
	var groupsWithUsers []GroupWithUsers
	err := json.NewDecoder(buf).Decode(&groupsWithUsers)
	if err != nil {
		return nil, err
	}
	return groupsWithUsers, nil
}

func GetLoadedJoinRequestFromRecorder(w *httptest.ResponseRecorder) (LoadedJoinRequest, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return LoadedJoinRequest{}, errors.New(buf.String())
	}

	var loadedJoinRequest LoadedJoinRequest
	err := json.NewDecoder(buf).Decode(&loadedJoinRequest)
	if err != nil {
		return LoadedJoinRequest{}, err
	}
	return loadedJoinRequest, nil
}

func GetLoadedJoinRequestsFromRecorder(w *httptest.ResponseRecorder) ([]LoadedJoinRequest, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return nil, errors.New(buf.String())
	}

	var loadedJoinRequests []LoadedJoinRequest
	err := json.NewDecoder(buf).Decode(&loadedJoinRequests)
	if err != nil {
		return nil, err
	}
	return loadedJoinRequests, nil
}

func GetJoinRequestFromRecorder(w *httptest.ResponseRecorder) (JoinRequest, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return JoinRequest{}, errors.New(buf.String())
	}

	var joinRequest JoinRequest
	err := json.NewDecoder(buf).Decode(&joinRequest)
	if err != nil {
		return JoinRequest{}, err
	}
	return joinRequest, nil
}

func GetJoinRequestsFromRecorder(w *httptest.ResponseRecorder) ([]JoinRequest, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return nil, errors.New(buf.String())
	}

	var joinRequests []JoinRequest
	err := json.NewDecoder(buf).Decode(&joinRequests)
	if err != nil {
		return nil, err
	}
	return joinRequests, nil
}

func GetJoinRequestRespondFromRecorder(w *httptest.ResponseRecorder) (JoinRequestRespond, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return JoinRequestRespond{}, errors.New(buf.String())
	}

	var joinRequestRespond JoinRequestRespond
	err := json.NewDecoder(buf).Decode(&joinRequestRespond)
	if err != nil {
		return JoinRequestRespond{}, err
	}
	return joinRequestRespond, nil
}

func GetMatchSettingFromRecorder(w *httptest.ResponseRecorder) (MatchSetting, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return MatchSetting{}, errors.New(buf.String())
	}

	var matchSetting MatchSetting
	err := json.NewDecoder(buf).Decode(&matchSetting)
	if err != nil {
		return MatchSetting{}, err
	}
	return matchSetting, nil
}

func GetMatchRequestFromRecorder(w *httptest.ResponseRecorder) (MatchRequest, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return MatchRequest{}, errors.New(buf.String())
	}

	var matchRequest MatchRequest
	err := json.NewDecoder(buf).Decode(&matchRequest)
	if err != nil {
		return MatchRequest{}, err
	}
	return matchRequest, nil
}

func GetLoadedMatchRequestFromRecorder(w *httptest.ResponseRecorder) (LoadedMatchRequest, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return LoadedMatchRequest{}, errors.New(buf.String())
	}

	var loadedMatchRequest LoadedMatchRequest
	err := json.NewDecoder(buf).Decode(&loadedMatchRequest)
	if err != nil {
		return LoadedMatchRequest{}, err
	}
	return loadedMatchRequest, nil
}

func GetIOReaderFromUser(user User) (io.Reader, error) {
	j, err := json.Marshal(user)
	if err != nil { return nil, err }
	return bytes.NewReader(j), nil
}


func GetIOReaderFromGroup(group Group) (io.Reader, error) {
	j, err := json.Marshal(group)
	if err != nil { return nil, err }
	return bytes.NewReader(j), nil
}

func GetIOReaderFromJoinRequestRespond(respond JoinRequestRespond) (io.Reader, error) {
	j, err := json.Marshal(respond)
	if err != nil { return nil, err }
	return bytes.NewReader(j), nil
}

func GetIOReaderFromJoinRequest(joinRequest JoinRequest) (io.Reader, error) {
	j, err := json.Marshal(joinRequest)
	if err != nil { return nil, err }
	return bytes.NewReader(j), nil
}

func GetIOReaderFromMatchSetting(matchSetting MatchSetting) (io.Reader, error) {
	j, err := json.Marshal(matchSetting)
	if err != nil { return nil, err }
	return bytes.NewReader(j), nil
}