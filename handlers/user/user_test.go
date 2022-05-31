package user;

import (
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

// Full test

func TestUserHandlers(t *testing.T) {
	t.Run("GetAllUsersHandler when DB is empty", testGetAllUsersHandlerWhenDBIsEmpty)
	t.Run("GetUserHandler when DB is empty", testGetUserHandlerWhenDBIsEmpty)
	t.Run("AddUserHandler", testAddUserHandler)
	t.Run("GetUserHandler", testGetUserHandler)
	t.Run("AddUserHandler no first name", testAddUserHandlerNoFirstName)
	t.Run("AddUserHandler no last name", testAddUserHandlerNoLastName)
	t.Run("AddUserHandler no gender", testAddUserHandlerNoGender)
	t.Run("AddUserHandler no faculty", testAddUserHandlerNoFaculty)
	t.Run("AddUserHandler no email", testAddUserHandlerNoEmail)
	t.Run("AddUserHandler no user role", testAddUserHandlerNoUserRole)
	t.Run("AddUserHandler same user", testAddSameUserHandler)
	t.Run("UpdateUserHandler unauthorized", testUpdateUserHandlerUnauthorized)
	t.Run("UpdateUserHandler authorized", testUpdateUserHandlerAuthorized)
	t.Run("DeleteUserHandler unauthorized", testDeleteUserHandlerUnauthorized)
	t.Run("DeleteUserHandler authorized", testDeleteUserHandlerAuthorised)
}

// Helpers

func getBufferFromRecorder(w *httptest.ResponseRecorder) *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Result().Body)
	return buf
}

func getUserFromRecorder(w *httptest.ResponseRecorder) (User, error) {
	buf := getBufferFromRecorder(w)
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

func getUsersFromRecorder(w *httptest.ResponseRecorder) ([]User, error) {
	buf := getBufferFromRecorder(w)
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

func getCookieFromRecorder(w *httptest.ResponseRecorder, name string) string {
	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}

func getIDCookieFromRecorder(w *httptest.ResponseRecorder) (string, int64, error) {
	sid := getCookieFromRecorder(w, "id")
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

func testGetAllUsersHandlerWhenDBIsEmpty(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user", nil)
	w := simulateRequest(req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllUser failed with status code of %d", w.Code)
	}
	users, err := getUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting user slice from recorder. %v", err)
	}
	if len(users) != 0 {
		t.Errorf("%d users found despite table being cleared", len(users))
	}
}

func testGetUserHandlerWhenDBIsEmpty(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user/1", nil)
	w := simulateRequest(req)
	if w.Code != http.StatusNotFound {
		t.Errorf("HTTP Request to GetUser did not have a status code of 404 not found")
	}
	_, err := getUserFromRecorder(w)
	if err == nil {
		t.Errorf("No error when getting a user that did not exist. %v", err)
	}
}

func testAddUserHandler(t *testing.T) {
	ioReaderUser, _ := getIOReaderFromUser(validUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to AddUser failed with status code of %d", w.Code)
	}
	var err error
	addedUser, err = getUserFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting addedUser from body. %v", err)
	}
	if addedUser.ID == 0 {
		t.Errorf("addedUser ID was not written by addUser call")
	}
	_, cookieID, err := getIDCookieFromRecorder(w)
	if err != nil || cookieID != addedUser.ID {
		t.Errorf("Error when retrieving id of cookie or id of cookie does not matched added User. %v", err)
	}
}

func testGetUserHandler(t *testing.T) {
	route := fmt.Sprintf("/user/%d", addedUser.ID)
	req, _ := http.NewRequest("GET", route, nil)
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetUser did not have a status code of 404 not found")
	}
	retrivedUser, err := getUserFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting user of id = %d from body. %v", addedUser.ID, err)
	}
	if !equal(retrivedUser, addedUser) {
		t.Errorf("retrieved user is not the same as the added user")
	}
}

func testAddUserHandlerNoFirstName(t *testing.T) {
	newUser := User{
		FirstName: "",
		LastName: validUser.LastName,
		Gender: validUser.Gender,
		Faculty: validUser.Faculty,
		Email: validUser.Email,
		UserRole: validUser.UserRole,
		Password: validUser.Password,
	}
	ioReaderUser, _ := getIOReaderFromUser(newUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := simulateRequest(req)
	if w.Code == http.StatusOK {
		t.Errorf("User with no first_name successfully added. Status Code: %d", w.Code)
	}
	matched, _ := regexp.MatchString("first_name", getBufferFromRecorder(w).String())
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of first_name")
	}
}

func testAddUserHandlerNoLastName(t *testing.T) {
	newUser := User{
		FirstName: validUser.FirstName,
		LastName: "",
		Gender: validUser.Gender,
		Faculty: validUser.Faculty,
		Email: validUser.Email,
		UserRole: validUser.UserRole,
		Password: validUser.Password,
	}
	ioReaderUser, _ := getIOReaderFromUser(newUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := simulateRequest(req)
	if w.Code == http.StatusOK {
		t.Errorf("User with no last_name successfully added. Status Code: %d", w.Code)
	}
	matched, _ := regexp.MatchString("last_name", getBufferFromRecorder(w).String())
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of last_name")
	}
}

func testAddUserHandlerNoGender(t *testing.T) {
	newUser := User{
		FirstName: validUser.FirstName,
		LastName: validUser.LastName,
		Gender: "",
		Faculty: validUser.Faculty,
		Email: validUser.Email,
		UserRole: validUser.UserRole,
		Password: validUser.Password,
	}
	ioReaderUser, _ := getIOReaderFromUser(newUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := simulateRequest(req)
	if w.Code == http.StatusOK {
		t.Errorf("User with no gender successfully added. Status Code: %d", w.Code)
	}
	matched, _ := regexp.MatchString("gender", getBufferFromRecorder(w).String())
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of gender")
	}
}

func testAddUserHandlerNoFaculty(t *testing.T) {
	newUser := User{
		FirstName: validUser.FirstName,
		LastName: validUser.LastName,
		Gender: validUser.Gender,
		Faculty: "",
		Email: validUser.Email,
		UserRole: validUser.UserRole,
		Password: validUser.Password,
	}
	ioReaderUser, _ := getIOReaderFromUser(newUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := simulateRequest(req)
	if w.Code == http.StatusOK {
		t.Errorf("User with no faculty successfully added. Status Code: %d", w.Code)
	}
	matched, _ := regexp.MatchString("faculty", getBufferFromRecorder(w).String())
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of faculty")
	}
}

func testAddUserHandlerNoEmail(t *testing.T) {
	newUser := User{
		FirstName: validUser.FirstName,
		LastName: validUser.LastName,
		Gender: validUser.Gender,
		Faculty: validUser.Faculty,
		Email: "",
		UserRole: validUser.UserRole,
		Password: validUser.Password,
	}
	
	ioReaderUser, _ := getIOReaderFromUser(newUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := simulateRequest(req)
	if w.Code == http.StatusOK {
		t.Errorf("User with no email successfully added. Status Code: %d", w.Code)
	}
	matched, _ := regexp.MatchString("email", getBufferFromRecorder(w).String())
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of email")
	}
}

func testAddUserHandlerNoUserRole(t *testing.T) {
	newUser := User{
		FirstName: validUser.FirstName,
		LastName: validUser.LastName,
		Gender: validUser.Gender,
		Faculty: validUser.Faculty,
		Email: validUser.Email,
		UserRole: "",
		Password: validUser.Password,
	}
	ioReaderUser, _ := getIOReaderFromUser(newUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := simulateRequest(req)
	if w.Code == http.StatusOK {
		t.Errorf("User with no user_role successfully added. Status Code: %d", w.Code)
	}
	matched, _ := regexp.MatchString("user_role", getBufferFromRecorder(w).String())
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of user_role")
	}
}

func testAddSameUserHandler(t *testing.T) {
	ioReaderUser, _ := getIOReaderFromUser(validUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := simulateRequest(req)
	if w.Code == http.StatusOK {
		t.Errorf("User with same details as addedUser successfully added. Status Code: %d", w.Code)
	}
}

func testUpdateUserHandlerUnauthorized(t *testing.T) {
	ioReaderUser, _ := getIOReaderFromUser(User{ FirstName: "UpdatedFirstName" })
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/user/%d", addedUser.ID), ioReaderUser)
	w := simulateRequest(req)
	_, err := getUserFromRecorder(w)
	matched, _ := regexp.MatchString(UnauthorizedErrorMessage, err.Error())
	if !matched {
		t.Errorf("Unauthorized user was able to make updates to user. %v", err)
	}
}

func testUpdateUserHandlerAuthorized(t *testing.T) {
	ioReaderUser, _ := getIOReaderFromUser(User{ FirstName: "UpdatedFirstName" })
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/user/%d", addedUser.ID), ioReaderUser)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", addedUser.ID),
	})
	w := simulateRequest(req)
	_, err := getUserFromRecorder(w)
	if err != nil {
		t.Errorf("Unable to update user despite being authorized. %v", err)
	}
}

func testDeleteUserHandlerUnauthorized(t *testing.T) {
	req, _ :=  http.NewRequest("DELETE", fmt.Sprintf("/user/%d", addedUser.ID), nil)
	w := simulateRequest(req)
	_, err := getUserFromRecorder(w)
	matched, _ := regexp.MatchString(UnauthorizedErrorMessage, err.Error())
	if !matched {
		t.Errorf("Unauthorized user was able to delete the user. %v", err)
	}
}

func testDeleteUserHandlerAuthorised(t *testing.T) {
	req, _ :=  http.NewRequest("DELETE", fmt.Sprintf("/user/%d", addedUser.ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", addedUser.ID),
	})
	w := simulateRequest(req)
	_, err := getUserFromRecorder(w)
	if err != nil {
		t.Errorf("Unable to delete user despite being authorized. %v", err)
	}
}