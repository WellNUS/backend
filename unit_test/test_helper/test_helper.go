package test_helper

import (
	"wellnus/backend/db/model"

	"database/sql"
	"regexp"
	"net/http"
	"net/http/httptest"
	"bytes"
	"errors"
	"encoding/json"
	"io"
	"fmt"
	"time"
	"math/rand"

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
type CounselRequest = model.CounselRequest
type UserIDBody = model.UserIDBody
type Event = model.Event
type EventWithUsers = model.EventWithUsers
type ProviderSetting = model.ProviderSetting
type Provider = model.Provider
type ProviderWithEvents = model.ProviderWithEvents
type Booking = model.Booking
type BookingUser = model.BookingUser
type BookingProvider = model.BookingProvider
type BookingRespond = model.BookingRespond

var ref_user_role 	[]string = []string{"MEMBER", "VOLUNTEER", "COUNSELLOR"}
var ref_category  	[]string = []string{"COUNSEL", "SUPPORT", "CUSTOM"}
var ref_faculty 	[]string = []string{"MIX", "SAME", "NONE"}
var ref_hobbies 	[]string = []string{"GAMING", "SINGING", "DANCING", "MUSIC", "SPORTS", "OUTDOOR", "BOOK", "ANIME", "MOVIES", "TV", "ART", "STUDY"}
var ref_mbti 		[]string = []string{"ISTJ","ISFJ","INFJ","INTJ","ISTP","ISFP","INFP","INTP","ESTP","ESFP","ENFP","ENTP","ESTJ","ESFJ","ENFJ","ENTJ"}
var ref_topics 		[]string = []string{"Anxiety", "OffMyChest", "SelfHarm"}
var ref_access 		[]string = []string{"PUBLIC", "PRIVATE"}

func ResetDB(db *sql.DB) {
	db.Exec("DELETE FROM wn_group")
	db.Exec("DELETE FROM wn_event")
	db.Exec("DELETE FROM wn_user")
}

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

func GetInt64FromRecorder(w *httptest.ResponseRecorder) (int64, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return 0, errors.New(buf.String())
	}

	var num int64
	err := json.NewDecoder(buf).Decode(&num)
	if err != nil {
		return 0, err
	}
	return num, nil
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

func GetCounselRequestFromRecorder(w *httptest.ResponseRecorder) (CounselRequest, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return CounselRequest{}, errors.New(buf.String())
	}
	var counselRequest CounselRequest
	err := json.NewDecoder(buf).Decode(&counselRequest)
	if err != nil {
		return CounselRequest{}, err
	}
	return counselRequest, nil
}

func GetCounselRequestsFromRecorder(w *httptest.ResponseRecorder) ([]CounselRequest, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return nil, errors.New(buf.String())
	}
	var counselRequests []CounselRequest
	err := json.NewDecoder(buf).Decode(&counselRequests)
	if err != nil {
		return nil, err
	}
	return counselRequests, nil
}

func GetEventFromRecorder(w *httptest.ResponseRecorder) (Event, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return Event{}, errors.New(buf.String())
	}
	var event Event
	err := json.NewDecoder(buf).Decode(&event)
	if err != nil {
		return Event{}, err
	}
	return event, nil
}

func GetEventsFromRecorder(w *httptest.ResponseRecorder) ([]Event, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return nil, errors.New(buf.String())
	}
	var events []Event
	err := json.NewDecoder(buf).Decode(&events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func GetEventWithUsersFromRecorder(w *httptest.ResponseRecorder) (EventWithUsers, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return EventWithUsers{}, errors.New(buf.String())
	}
	var eventWithUsers EventWithUsers
	err := json.NewDecoder(buf).Decode(&eventWithUsers)
	if err != nil {
		return EventWithUsers{}, err
	}
	return eventWithUsers, nil
}

func GetEventsWithUsersFromRecorder(w *httptest.ResponseRecorder) ([]EventWithUsers, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return nil, errors.New(buf.String())
	}
	var eventsWithUsers []EventWithUsers
	err := json.NewDecoder(buf).Decode(&eventsWithUsers)
	if err != nil {
		return nil, err
	}
	return eventsWithUsers, nil
}

func GetProviderSettingFromRecorder(w *httptest.ResponseRecorder) (ProviderSetting, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return ProviderSetting{}, errors.New(buf.String())
	}
	var providerSetting ProviderSetting
	err := json.NewDecoder(buf).Decode(&providerSetting)
	if err != nil {
		return ProviderSetting{}, err
	}
	return providerSetting, nil
}

func GetProviderSettingsFromRecorder(w *httptest.ResponseRecorder) ([]ProviderSetting, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return nil, errors.New(buf.String())
	}
	var providerSettings []ProviderSetting
	err := json.NewDecoder(buf).Decode(&providerSettings)
	if err != nil {
		return nil, err
	}
	return providerSettings, nil
}

func GetProviderFromRecorder(w *httptest.ResponseRecorder) (Provider, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return Provider{}, errors.New(buf.String())
	}
	var provider Provider
	err := json.NewDecoder(buf).Decode(&provider)
	if err != nil {
		return Provider{}, err
	}
	return provider, nil
}

func GetProvidersFromRecorder(w *httptest.ResponseRecorder) ([]Provider, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return nil, errors.New(buf.String())
	}
	var providers []Provider
	err := json.NewDecoder(buf).Decode(&providers)
	if err != nil {
		return nil, err
	}
	return providers, nil
}

func GetProviderWithEventsFromRecorder(w *httptest.ResponseRecorder) (ProviderWithEvents, error) {
	buf := GetBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return ProviderWithEvents{}, errors.New(buf.String())
	}
	var providerWithEvents ProviderWithEvents
	err := json.NewDecoder(buf).Decode(&providerWithEvents)
	if err != nil {
		return ProviderWithEvents{}, err
	}
	return providerWithEvents, nil
}

func CheckErrorMessageFromRecorder(w *httptest.ResponseRecorder, pattern string) (string, bool) {
	errString := GetBufferFromRecorder(w).String()
	matched, _ := regexp.MatchString(pattern, errString)
	return errString, matched
}

func GetIOReaderFromObject(obj interface{}) (io.Reader, error) {
	j, err := json.Marshal(obj)
	if err != nil { return nil, err }
	return bytes.NewReader(j), nil
}

func GenerateRandomString(l int) string {
	CharSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	Rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, l)
	charSetLen := len(CharSet)
	for i := range b {
		b[i] = CharSet[Rand.Intn(charSetLen)]
	}
	return string(b)
}

func GetTestUser(i int) User {
	email := GenerateRandomString(20)
	role := ref_user_role[i % len(ref_user_role)]

	return User{
		FirstName: fmt.Sprintf("TestUser%d", i),
		LastName: fmt.Sprintf("TestLastName%d", i),
		Gender: "M",
		Faculty: "COMPUTING",
		Email: fmt.Sprintf("%s@u.nus.edu", email),
		UserRole: role,
		Password: "123",
		PasswordHash: "",
	}
}

func GetTestGroup(i int) Group {
	category := ref_category[i % len(ref_category)]
	return Group{
		GroupName: fmt.Sprintf("NewGroupName%d", i),
		GroupDescription: "NewGroupDescription",
		Category: category,
	}
}

func GetRandomTestMatchSetting() MatchSetting {
	Rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	facultyPreference := ref_faculty[Rand.Intn(len(ref_faculty))]
	mbti := ref_mbti[Rand.Intn(len(ref_mbti))]
	hobbies := make([]string, 0)
	for _, hobby := range ref_hobbies {
		if Rand.Intn(3) == 1 { hobbies = append(hobbies, hobby) }
		if len(hobbies) >= 4 { break }
	}
	matchSetting := MatchSetting{
		FacultyPreference: facultyPreference,
		Hobbies: hobbies,
		MBTI: mbti,
	}
	return matchSetting
}

func GetTestCounselRequest(i int) CounselRequest {
	counselRequest := CounselRequest{
		Nickname: "testRecipient",
		Details: "I am stressed",
		Topics: []string{ref_topics[0], ref_topics[1]}}
	if i % 3 == 1 {
		counselRequest.Topics = []string{ref_topics[1], ref_topics[2]}
	} else if i % 3 == 2 {
		counselRequest.Topics = []string{ref_topics[0], ref_topics[2]}
	}
	return counselRequest
}

func GetTestEvent(i int) Event {
	startTime, _ := time.Parse(time.RFC3339, "2050-01-01T08:00:00Z08:00")
	endTime, _ := time.Parse(time.RFC3339, "2055-01-01T08:00:00Z08:00")
	access := ref_access[i % len(ref_access)]
	category := ref_category[i % len(ref_category)]
	return Event{
		EventName: fmt.Sprintf("TestEvent%d", i),
		EventDescription: "NewEventDescription",
		StartTime: startTime,
		EndTime: endTime,
		Access: access,
		Category: category,
	}
}

func GetTestProviderSetting(i int) ProviderSetting {
	providerSetting := ProviderSetting{
		Intro: "I am a professional counsellor",
		Topics: []string{ref_topics[0], ref_topics[1]}}
	if i % 3 == 1 {
		providerSetting.Topics = []string{ref_topics[1], ref_topics[2]}
	} else if i % 3 == 2 {
		providerSetting.Topics = []string{ref_topics[0], ref_topics[2]}
	}
	return providerSetting
}

func GetTestBooking(i int, recipientID int64, providerID int64) Booking {
	startTime, _ := time.Parse(time.RFC3339, "2050-01-01T08:00:00Z08:00")
	endTime, _ := time.Parse(time.RFC3339, "2055-01-01T08:00:00Z08:00")
	return Booking{
		ProviderID: providerID,
		Nickname: fmt.Sprintf("TestNickName%d", i),
		Details: "Looking to talk about my difficulties",
		StartTime: startTime,
		EndTime: endTime,
	}
}

func SetupUsers(db *sql.DB, num int) ([]User, error) {
	users := make([]User, num)
	for i := 0; i < num; i++ {
		user, err := model.AddUser(db, GetTestUser(i))
		if err != nil { return nil, err }
		users[i] = user
	}
	return users, nil
}

func SetupGroupsForUsers(db *sql.DB, users []User) ([]Group, error) {
	groups := make([]Group, len(users))
	for i, user := range users {
		groupWithUsers, err := model.AddGroupWithUserIDs(db, GetTestGroup(i), []int64{ user.ID })
		if err != nil { return nil, err }
		groups[i] = groupWithUsers.Group
	}
	return groups, nil
}

func SetupMatchSettingForUsers(db *sql.DB, users []User) ([]MatchSetting, error) {
	matchSettings := make([]MatchSetting, len(users))
	for i, user := range users {
		matchSetting, err := model.AddUpdateMatchSettingOfUser(db, GetRandomTestMatchSetting(), user.ID)
		if err != nil { return nil, err }
		matchSettings[i] = matchSetting
	}
	return matchSettings, nil
}

func SetupSessionForUsers(db *sql.DB, users []User) ([]string, error) {
	sessionKeys := make([]string, len(users))
	for i, user := range users {
		sessionKey, err := model.CreateNewSession(db, user.ID)
		if err != nil { return nil, err }
		sessionKeys[i] = sessionKey
	}
	return sessionKeys, nil
}

func SetupMatchRequestForUsers(db *sql.DB, users []User) ([]MatchRequest, error) {
	matchRequests := make([]MatchRequest, len(users))
	for i, user := range users {
		matchRequest, err := model.AddMatchRequest(db, user.ID)
		if err != nil { return nil, err }
		matchRequests[i] = matchRequest
	}
	return matchRequests, nil
}

func SetupCounselRequestForUsers(db *sql.DB, users []User) ([]CounselRequest, error) {
	counselRequests := make([]CounselRequest, len(users))
	for i, user := range users {
		counselRequest, err := model.AddUpdateCounselRequest(db, GetTestCounselRequest(i), user.ID)
		if err != nil { return nil, err }
		counselRequests[i] = counselRequest
	}
	return counselRequests, nil
}

func SetupEventForUsers(db *sql.DB, users []User) ([]Event, error) {
	events := make([]Event, len(users))
	for i, user := range users {
		eventWithUsers, err := model.AddEventWithUserIDs(db, GetTestEvent(i), []int64{user.ID})
		if err != nil { return nil, err }
		events[i] = eventWithUsers.Event
	}
	return events, nil
}

func SetupProviderSettingForUsers(db *sql.DB, users []User) ([]ProviderSetting, error) {
	providerSettings := make([]ProviderSetting, len(users))
	for i, user := range users {
		providerSetting, err := model.AddUpdateProviderSettingOfUser(db, GetTestProviderSetting(i), user.ID)
		if err != nil { return nil, err }
		providerSettings[i] = providerSetting
	}
	return providerSettings, nil
}