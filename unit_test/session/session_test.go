package session

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper/http_error"
	"wellnus/backend/unit_test/test_helper"

	"testing"
	"net/http"
)

// Full Tests

func TestSession(t *testing.T) {
	t.Run("Successful Login Handler", testSuccessfulLoginHandler)
	t.Run("Failed Login Handler", testFailedLoginHandler)
	t.Run("Logout Handler", testLogoutHandler)
}

// Helpers

func testSuccessfulLoginHandler(t *testing.T) {
	password := test_helper.GetTestUser(0).Password
	loginAttempt := User{
		Email: testUsers[0].Email, 
		Password: password}
	IOReaderAttempt, _ := test_helper.GetIOReaderFromObject(loginAttempt)
	req, _ := http.NewRequest("POST", "/session", IOReaderAttempt)
	w := test_helper.SimulateRequest(Router, req)
	sessionResponse, err := test_helper.GetSessionResponseFromRecorder(w)
	if err != nil { t.Errorf("An error occured while retrieving response body. %v", err)}
	if !sessionResponse.LoggedIn { t.Errorf("Not logged in despite logging in") }
	sessionKey = test_helper.GetCookieFromRecorder(w, "session_key")
	userID, err := model.GetUserIDFromSessionKey(DB, sessionKey)
	if err != nil { t.Errorf("An error occured while retrieving userID from session key. %v", err)}
	if userID != sessionResponse.User.ID { 
		t.Errorf("Logged in as a user of id = %d instead of correct user of id = %d", userID, sessionResponse.User.ID)
	}
}

func testFailedLoginHandler(t *testing.T) {
	loginAttempt := User{
		Email: testUsers[0].Email, 
		Password: "WrongPassword"}
	IOReaderAttempt, _ := test_helper.GetIOReaderFromObject(loginAttempt)
	req, _ := http.NewRequest("POST", "/session", IOReaderAttempt)
	w := test_helper.SimulateRequest(Router, req)
	sessionResponse, err := test_helper.GetSessionResponseFromRecorder(w)
	if err != nil { t.Errorf("An error occured while retrieving response body. %v", err)}
	if sessionResponse.LoggedIn { t.Errorf("Logged in despite wrong password") }
}

func testLogoutHandler(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/session", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKey,
	})
	w := test_helper.SimulateRequest(Router, req)
	sessionResponse, err := test_helper.GetSessionResponseFromRecorder(w)
	if err != nil { t.Errorf("An error occured while retrieving response body. %v", err)}
	if sessionResponse.LoggedIn { t.Errorf("response indicate that logout was unsuccessful") }
	newSessionKey := test_helper.GetCookieFromRecorder(w, "session_key")
	if newSessionKey != "" { t.Errorf("Session Key cookie is still present after logout. SessionKey = '%s'", newSessionKey) }

	//Check if session is still stored in DB
	_, err = model.GetUserIDFromSessionKey(DB, sessionKey)
	if err != http_error.UnauthorizedError {
		t.Errorf("Session still exist in DB as no unauthorized error was thrown. %v", err)
	}
}