package user

import (
	"testing"
	"regexp"
)

// Main test

func TestUserQuery(t *testing.T) {
	t.Run("GetAllUsers when DB is empty", testGetAllUsersWhenDBIsEmpty)
	t.Run("GetUser when DB is empty", testGetUserWhenDBIsEmpty)
	t.Run("AddUser", testAddUser)
	t.Run("GetUser", testGetUser)
	t.Run("AddUser no first name", testAddUserNoFirstName)
	t.Run("AddUser no last name", testAddUserNoLastName)
	t.Run("AddUser no gender", testAddUserNoGender)
	t.Run("AddUser no faculty", testAddUserNoFaculty)
	t.Run("AddUser no email", testAddUserNoEmail)
	t.Run("AddUser no user role", testAddUserNoUserRole)
	t.Run("AddUser same user", testAddSameUser)
	t.Run("UpdateUser", testUpdateUser)
	t.Run("GetUser after update", testGetUserAfterUpdate)
	t.Run("DeleteUser unauthorized", testDeleteUser)
	t.Run("GetUser after delete", testGetUserAfterDelete)
}

// Helpers

func testGetAllUsersWhenDBIsEmpty(t *testing.T) {
	users, err := GetAllUsers(db)
	if err != nil {
		t.Errorf("Error when getting all users. %v", err)
	}
	if len(users) != 0 {
		t.Errorf("%d users found despite table being cleared", len(users))
	}
}

func testGetUserWhenDBIsEmpty(t *testing.T) {
	_, err := GetUser(db, 1)
	if err == nil || err.Error() != "404 Not Found" {
		t.Errorf("Expected a not found error but either got no error or a different error. %v", err)
	}
}

func testAddUser(t *testing.T) {
	var err error
	addedUser, err = AddUser(db, validUser)
	if err != nil { 
		t.Errorf("An error occured while adding a new user. %v", err) 
	}
	if addedUser.ID == 0 {
		t.Errorf("addedUser ID not written by AddUser call")
	}
}

func testGetUser(t *testing.T) {
	// Checking of added User
	retrievedUser, err := GetUser(db, addedUser.ID)
	if err != nil { 
		t.Errorf("An error occured while retrieving user of id = %d. %v", addedUser.ID, err)
	}
	if !equal(addedUser, retrievedUser) {
		t.Errorf("retrieved user is not the same as the added user")
	}
}

func testAddUserNoFirstName(t *testing.T) {
	newUser := User{
		FirstName: "",
		LastName: validUser.LastName,
		Gender: validUser.Gender,
		Faculty: validUser.Faculty,
		Email: validUser.Email,
		UserRole: validUser.UserRole,
		Password: validUser.Password,
	}
	_, err := AddUser(db, newUser)
	if err == nil {
		t.Errorf("User without first name was successfully added")
	}
	matched, _ := regexp.MatchString("first_name", err.Error())
	if !matched {
		t.Errorf("Error did not contain any instance of first_name. %v", err)
	}
}

func testAddUserNoLastName(t *testing.T) {
	newUser := User{
		FirstName: validUser.FirstName,
		LastName: "",
		Gender: validUser.Gender,
		Faculty: validUser.Faculty,
		Email: validUser.Email,
		UserRole: validUser.UserRole,
		Password: validUser.Password,
	}
	_, err := AddUser(db, newUser)
	if err == nil {
		t.Errorf("User without last name was successfully added")
	}
	matched, _ := regexp.MatchString("last_name", err.Error())
	if !matched {
		t.Errorf("Error did not contain any instance of last_name. %v", err)
	}
}

func testAddUserNoGender(t *testing.T) {
	newUser := User{
		FirstName: validUser.FirstName,
		LastName: validUser.LastName,
		Gender: "",
		Faculty: validUser.Faculty,
		Email: validUser.Email,
		UserRole: validUser.UserRole,
		Password: validUser.Password,
	}
	_, err := AddUser(db, newUser)
	if err == nil {
		t.Errorf("User without gender was successfully added")
	}
	matched, _ := regexp.MatchString("gender", err.Error())
	if !matched {
		t.Errorf("Error did not contain any instance of gender. %v", err)
	}
}

func testAddUserNoFaculty(t *testing.T) {
	newUser := User{
		FirstName: validUser.FirstName,
		LastName: validUser.LastName,
		Gender: validUser.Gender,
		Faculty: "",
		Email: validUser.Email,
		UserRole: validUser.UserRole,
		Password: validUser.Password,
	}
	_, err := AddUser(db, newUser)
	if err == nil {
		t.Errorf("User without faculty was successfully added")
	}
	matched, _ := regexp.MatchString("faculty", err.Error())
	if !matched {
		t.Errorf("Error did not contain any instance of faculty. %v", err)
	}
}

func testAddUserNoEmail(t *testing.T) {
	newUser := User{
		FirstName: validUser.FirstName,
		LastName: validUser.LastName,
		Gender: validUser.Gender,
		Faculty: validUser.Faculty,
		Email: "",
		UserRole: validUser.UserRole,
		Password: validUser.Password,
	}
	_, err := AddUser(db, newUser)
	if err == nil {
		t.Errorf("User without email was successfully added")
	}
	matched, _ := regexp.MatchString("email", err.Error())
	if !matched {
		t.Errorf("Error did not contain any instance of email. %v", err)
	}
}

func testAddUserNoUserRole(t *testing.T) {
	newUser := User{
		FirstName: validUser.FirstName,
		LastName: validUser.LastName,
		Gender: validUser.Gender,
		Faculty: validUser.Faculty,
		Email: validUser.Email,
		UserRole: "",
		Password: validUser.Password,
	}
	_, err := AddUser(db, newUser)
	if err == nil {
		t.Errorf("User without user_role was successfully added")
	}
	matched, _ := regexp.MatchString("user_role", err.Error())
	if !matched {
		t.Errorf("Error did not contain any instance of user_role. %v", err)
	}
}

func testAddSameUser(t *testing.T) {
	_, err := AddUser(db, validUser)
	if err == nil { 
		t.Errorf("User with same already existing email was added") 
	}
}

func testUpdateUser(t *testing.T) {
	var err error
	newFirstName := "UpdatedFirstName"
	addedUser, err = UpdateUser(db, User{ FirstName: newFirstName }, addedUser.ID)
	if err != nil {
		t.Errorf("An error occured while updating user. %v", err)
	}
	if addedUser.FirstName != newFirstName {
		t.Errorf("Returned updated object did not reflect updates")
	}
}

func testGetUserAfterUpdate(t *testing.T) {
	retrievedUser, err := GetUser(db, addedUser.ID)
	if err != nil { 
		t.Errorf("An error occured while retrieving user of id = %d. %v", addedUser.ID, err)
	}
	if !equal(addedUser, retrievedUser) {
		t.Errorf("Updates made to added user was not reflected in database")
	}
}

func testDeleteUser(t *testing.T) {
	user, err := DeleteUser(db, addedUser.ID)
	if err != nil {
		t.Errorf("An error occured while deleting user of id = %d. %v", addedUser.ID, err)
	}
	if user.ID != addedUser.ID {
		t.Errorf("Returned user id did not matched the original id of added user")
	}
}

func testGetUserAfterDelete(t *testing.T) {
	_, err := GetUser(db, addedUser.ID)
	if err == nil || err.Error() != "404 Not Found" {
		t.Errorf("Expected a not found error but either got no error or a different error. %v", err)
	}
}