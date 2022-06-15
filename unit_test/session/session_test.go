package session

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"io"
	"bytes"
	"errors"
	"encoding/json"
	"strconv"
)

// Full Tests

func TestSession(t *testing.T) {
	t.Run("Successful Login Handler", testSuccessfulLoginHandler)
	t.Run("Failed Login Handler", testFailedLoginHandler)
	t.Run("Logout Handler", testLogoutHandler)
}

// Helpers

func getResp(w *httptest.ResponseRecorder) (Resp, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Result().Body)
	if w.Code != http.StatusOK {
		return Resp{}, errors.New(buf.String())
	}

	//fmt.Printf("Response Body: %v \n", buf)
	var resp Resp
	err := json.NewDecoder(buf).Decode(&resp)
	if err != nil {
		return Resp{}, err
	}
	return resp, nil
}

func getCookie(w *httptest.ResponseRecorder, name string) string {
	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}

func getIDCookie(w *httptest.ResponseRecorder) (string, int64, error) {
	sid := getCookie(w, "id")
	id, err := strconv.ParseInt(sid, 0, 64)
	return sid, id, err
}

func simulateRequest(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	Router.ServeHTTP(w, req)
	return w
}

func getIOReaderFromUser(user User) (io.Reader, error) {
	j, err := json.Marshal(user)
	if err != nil { return nil, err }
	return bytes.NewReader(j), nil
}

func testSuccessfulLoginHandler(t *testing.T) {
	loginAttempt := User{
		Email: validUser.Email, 
		Password: validUser.Password}
	IOReaderAttempt, _ := getIOReaderFromUser(loginAttempt)
	req, _ := http.NewRequest("POST", "/session", IOReaderAttempt)
	w := simulateRequest(req)
	resp, err := getResp(w)
	if err != nil { t.Errorf("An error occured while retrieving response body. %v", err)}
	if !resp.LoggedIn { t.Errorf("Not logged in despite logging in") }
	_, id, err := getIDCookie(w)
	if err != nil { t.Errorf("An error occured while retrieving ID cookie. %v", err)}
	if id != resp.User.ID { t.Errorf("Logged in as a user of id = %d instead of correct user of id = %d", id, resp.User.ID) }
}

func testFailedLoginHandler(t *testing.T) {
	loginAttempt := User{
		Email: validUser.Email, 
		Password: "WrongPassword"}
	IOReaderAttempt, _ := getIOReaderFromUser(loginAttempt)
	req, _ := http.NewRequest("POST", "/session", IOReaderAttempt)
	w := simulateRequest(req)
	resp, err := getResp(w)
	if err != nil { t.Errorf("An error occured while retrieving response body. %v", err)}
	if resp.LoggedIn { t.Errorf("Logged in despite wrong password") }
}

func testLogoutHandler(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/session", nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: "9999",
	})
	w := simulateRequest(req)
	resp, err := getResp(w)
	if err != nil { t.Errorf("An error occured while retrieving response body. %v", err)}
	if resp.LoggedIn { t.Errorf("response indicate that logout was unsuccessful") }
	_, _, err = getIDCookie(w)
	if err == nil { t.Errorf("ID cookie is still present after logout") }
}