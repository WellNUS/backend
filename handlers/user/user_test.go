package user;

import (
	"wellnus/backend/handlers/httpError"

	"fmt"
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
	"strconv"
	"errors"
	"regexp"
	"io"

	"encoding/json"
)

var (
	NotFoundErrorMessage 		string = httpError.NotFoundError.Error()
	UnauthorizedErrorMessage	string = httpError.UnauthorizedError.Error()
)

var validUser User = User{
	FirstName: "NewFirstName",
	LastName: "NewLastName",
	Gender: "M",
	Faculty: "COMPUTING",
	Email: "NewEmail@u.nus.edu",
	UserRole: "VOLUNTEER",
	Password: "NewPassword",
	PasswordHash: "",
}

// Helper functions
func getUser(w *httptest.ResponseRecorder) (User, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Result().Body)
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
	router.ServeHTTP(w, req)
	return w
}

func getIOReaderFromUser(user User) (io.Reader, error) {
	j, err := json.Marshal(user)
	if err != nil { return nil, err }
	return bytes.NewReader(j), nil
}

// Main tests
func TestGetAllUsersHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user", nil)
	w := simulateRequest(req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllUser failed with status code of %d", w.Code)
	}
}

func TestGetUserHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user/1", nil)
	w := simulateRequest(req)
	user, err := getUser(w)
	if err != nil {
		t.Errorf("An error occured while getting User, %v", err)
	}
	if user.ID != 1 {
		t.Errorf("HTTP Request to GetUser of id = 1 did not get user of given ID")
	}
}

// Checks for error on requests and cookies. Does not check for if the database is properly updated
func TestAddUpdateAndRemoveUserHandler(t *testing.T) {
	// Adding new user
	ioReaderUser, _ := getIOReaderFromUser(validUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := simulateRequest(req)
	user, err := getUser(w)
	if err != nil {
		t.Errorf("Error occured while getting User, %v", err)
	}
	sCookieID, cookieID, err := getIDCookie(w)
	if err != nil || cookieID != user.ID {
		t.Errorf("Error when retrieving id of cookie or id of cookie does not matched added User. %v", err)
	}
	
	// Updating new user
	ioReaderUser, _ = getIOReaderFromUser(User{ FirstName: "UpdatedFirstName" })
	// Updating new user unauthorized
	req, _ = http.NewRequest("PATCH", fmt.Sprintf("/user/%d", user.ID), ioReaderUser)
	w = simulateRequest(req)
	_, err = getUser(w)
	matched, _ := regexp.MatchString(UnauthorizedErrorMessage, err.Error())
	if !matched {
		t.Errorf("Unauthorized user was able to make updates to user. %v", err)
	}

	// Updating new user authorized
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: sCookieID,
	})
	w = simulateRequest(req)
	_, err = getUser(w)
	if err != nil {
		t.Errorf("Unable to update user despite being authorized. %v", err)
	}

	// Removing new user unauthorized
	req, _ =  http.NewRequest("DELETE", fmt.Sprintf("/user/%d", user.ID), nil)
	w = simulateRequest(req)
	_, err = getUser(w)
	matched, _ = regexp.MatchString(UnauthorizedErrorMessage, err.Error())
	if !matched {
		t.Errorf("Unauthorized user was able to make updates to user. %v", err)
	}
	// Removing new user authorized
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: sCookieID,
	})
	w = simulateRequest(req)
	_, err = getUser(w)
	if err != nil {
		t.Errorf("Unable to delete user despite being authorized. %v", err)
	}
}