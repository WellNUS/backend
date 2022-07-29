package user

import (
	"wellnus/backend/db/model"
	"wellnus/backend/unit_test/test_helper"

	"fmt"
	"testing"
	"net/http"
)

// Full test

func TestUserHandlers(t *testing.T) {
	t.Run("GetAllUsersHandler", testGetAllUsersHandler)
	t.Run("GetAllUsersHandlerMember", testGetAllUsersHandlerMember)
	t.Run("GetAllUsersHandlerVolunteer", testGetAllUsersHandlerVolunteer)
	t.Run("GetAllUsersHandlerCounsellor", testGetAllUsersHandlerCounsellor)	
	t.Run("GetAllUsersHandlerProvider", testGetAllUsersHandlerProvider)
	t.Run("GetUserHandler miss", testGetUserHandlerMiss)
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

func testGetAllUsersHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllUser failed with status code of %d", w.Code)
	}
	users, err := test_helper.GetUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting user slice from recorder. %v", err)
	}
	if len(users) != 3 {
		t.Errorf("%d users found despite table setup with 3 users", len(users))
	}
}

func testGetAllUsersHandlerMember(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user?role=MEMBER", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllUser failed with status code of %d", w.Code)
	}
	users, err := test_helper.GetUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting user slice from recorder. %v", err)
	}
	if len(users) != 1 {
		t.Errorf("%d users found despite table setup with 1 member", len(users))
	}
}

func testGetAllUsersHandlerVolunteer(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user?role=VOLUNTEER", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllUser failed with status code of %d", w.Code)
	}
	users, err := test_helper.GetUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting user slice from recorder. %v", err)
	}
	if len(users) != 1 {
		t.Errorf("%d users found despite table setup with 1 volunteer", len(users))
	}
}

func testGetAllUsersHandlerCounsellor(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user?role=COUNSELLOR", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllUser failed with status code of %d", w.Code)
	}
	users, err := test_helper.GetUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting user slice from recorder. %v", err)
	}
	if len(users) != 1 {
		t.Errorf("%d users found despite table setup with 1 counsellor", len(users))
	}
}

func testGetAllUsersHandlerProvider(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user?role=PROVIDER", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllUser failed with status code of %d", w.Code)
	}
	users, err := test_helper.GetUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting user slice from recorder. %v", err)
	}
	if len(users) != 2 {
		t.Errorf("%d users found despite table setup with 2 provider", len(users))
	}
}

func testGetUserHandlerMiss(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user/999", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("HTTP Request to GetUser did not have a status code of 404 not found")
	}
	_, err := test_helper.GetUserFromRecorder(w)
	if err == nil {
		t.Errorf("No error when getting a user that did not exist. %v", err)
	}
}

func testAddUserHandler(t *testing.T) {
	ioReaderUser, _ := test_helper.GetIOReaderFromObject(validUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to AddUser failed with status code of %d", w.Code)
	}
	userWithSessionKey, err := test_helper.GetUserWithSessionKeyFromRecorder(w)
	addedUser = userWithSessionKey.User
	if err != nil {
		t.Errorf("An error occured while getting addedUser from body. %v", err)
	}
	if addedUser.ID == 0 {
		t.Errorf("addedUser ID was not written by addUser call")
	}
	sessionKey = userWithSessionKey.SessionKey
	userID, err := model.GetUserIDFromSessionKey(DB, sessionKey)
	if err != nil || userID != addedUser.ID {
		t.Errorf("Error when retrieving userID from sessionKey or the userID does not matched added User. %v", err)
	}
}

func testGetUserHandler(t *testing.T) {
	route := fmt.Sprintf("/user/%d", addedUser.ID)
	req, _ := http.NewRequest("GET", route, nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetUser did not have a status code of 404 not found")
	}
	retrivedUserWithGroups, err := test_helper.GetUserWithGroupsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting user of id = %d from body. %v", addedUser.ID, err)
	}
	if !retrivedUserWithGroups.User.Equal(addedUser) {
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
	ioReaderUser, _ := test_helper.GetIOReaderFromObject(newUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code == http.StatusOK {
		t.Errorf("User with no first_name successfully added. Status Code: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, "first_name")
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of first_name. %s", errString)
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
	ioReaderUser, _ := test_helper.GetIOReaderFromObject(newUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code == http.StatusOK {
		t.Errorf("User with no last_name successfully added. Status Code: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, "last_name")
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of last_name. %s", errString)
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
	ioReaderUser, _ := test_helper.GetIOReaderFromObject(newUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code == http.StatusOK {
		t.Errorf("User with no gender successfully added. Status Code: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, "gender")
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of gender. %s", errString)
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
	ioReaderUser, _ := test_helper.GetIOReaderFromObject(newUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code == http.StatusOK {
		t.Errorf("User with no faculty successfully added. Status Code: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, "faculty")
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of faculty. %s", errString)
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
	
	ioReaderUser, _ := test_helper.GetIOReaderFromObject(newUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code == http.StatusOK {
		t.Errorf("User with no email successfully added. Status Code: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, "email")
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of email. %s", errString)
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
	ioReaderUser, _ := test_helper.GetIOReaderFromObject(newUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code == http.StatusOK {
		t.Errorf("User with no user_role successfully added. Status Code: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, "user_role")
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of user_role. %s", errString)
	}
}

func testAddSameUserHandler(t *testing.T) {
	ioReaderUser, _ := test_helper.GetIOReaderFromObject(validUser)
	req, _ := http.NewRequest("POST", "/user", ioReaderUser)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code == http.StatusOK {
		t.Errorf("User with same details as addedUser successfully added. Status Code: %d", w.Code)
	}
}

func testUpdateUserHandlerUnauthorized(t *testing.T) {
	ioReaderUser, _ := test_helper.GetIOReaderFromObject(User{ FirstName: "UpdatedFirstName" })
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/user/%d", addedUser.ID), ioReaderUser)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized{
		t.Errorf("Unauthorised status code not given. Status Code: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, UnauthorizedErrorMessage)
	if !matched {
		t.Errorf("Unauthorized user was not unauthorised. %s", errString)
	}
}

func testUpdateUserHandlerAuthorized(t *testing.T) {
	ioReaderUser, _ := test_helper.GetIOReaderFromObject(User{ FirstName: "UpdatedFirstName" })
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/user/%d", addedUser.ID), ioReaderUser)
	req.Header.Add("session_key", sessionKey)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to updateUserHandler did not statusOk. Status Code: %d", w.Code)
	}
	_, err := test_helper.GetUserFromRecorder(w)
	if err != nil {
		t.Errorf("Unable to update user despite being authorized. %v", err)
	}
}

func testDeleteUserHandlerUnauthorized(t *testing.T) {
	req, _ :=  http.NewRequest("DELETE", fmt.Sprintf("/user/%d", addedUser.ID), nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Unauthorised status code not given. Status Code: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, UnauthorizedErrorMessage)
	if !matched {
		t.Errorf("Unauthorized user was able to delete the user. %s", errString)
	}
}

func testDeleteUserHandlerAuthorised(t *testing.T) {
	req, _ :=  http.NewRequest("DELETE", fmt.Sprintf("/user/%d", addedUser.ID), nil)
	req.Header.Add("session_key", sessionKey)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to updateUserHandler did not statusOk. Status Code: %d", w.Code)
	}
	_, err := test_helper.GetUserFromRecorder(w)
	if err != nil {
		t.Errorf("Unable to delete user despite being authorized. %v", err)
	}
}