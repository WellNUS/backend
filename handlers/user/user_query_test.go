package user

import (
	"testing"
	"errors"
	"regexp"
)

var templateUser User = User{
	FirstName: "NewFirstName",
	LastName: "NewLastName",
	Gender: "M",
	Email: "NewEmail@u.nus.edu",
	UserRole: "VOLUNTEER",
	Password: "NewPassword",
	PasswordHash: "",
}

func equal(user1 User, user2 User) bool {
	return user1.ID == user2.ID &&
	user1.FirstName == user2.FirstName &&
	user1.LastName == user2.LastName &&
	user1.Gender == user2.Gender &&
	user1.Email == user2.Email &&
	user1.UserRole == user2.UserRole &&
	user1.PasswordHash == user2.PasswordHash
}

func testAddPatchRemoveUser(newUser User) error {
	// Add new User to DB
	newUser, err := AddUser(db, newUser)
	if err != nil { return err }
	id := newUser.ID

	// Retrieve new User. Check if equal with original
	retrievedUser, err := GetUser(db, id)
	if err != nil { return err }
	if !equal(newUser, retrievedUser) { return errors.New("Storage of new user not done properly") }

	// Update newUser on DB
	_, err = UpdateUser(db, User{ FirstName: "UpdatedFirstName" }, id)
	if err != nil { return err }

	// Retrieve updated User. Check if equal with updated original
	newUser.FirstName = "UpdatedFirstName"
	retrievedUser, err = GetUser(db, id)
	if err != nil { return err }
	if !equal(newUser, retrievedUser) { return errors.New("Update of new user not done properly") }

	//Remove User from DB
	_, err = DeleteUser(db, id)
	if err != nil { return err }

	// Attempt to retrieve deleted 
	retrievedUser, err = GetUser(db, id)
	if err == nil { return errors.New("User was not deleted from database upon deletion") }
	if err.Error() != "404 Not Found" { return errors.New("Error was thrown but it was not '404 Not Found'") }
	return nil
}

func TestGetUser(t *testing.T) {
	user, err := GetUser(db, 999999)
	if err.Error() != "404 Not Found" {
		t.Errorf("Expected a not found error but got a different error. %v", err)
	}
	user, err = GetUser(db, 1)
	if err != nil {
		t.Errorf("Error when retrieving user of id = 1. %v", err)
	}
	if user.ID != 1 {
		t.Errorf("Expected retrived user to have and id of 1. But has id of %d", user.ID)
	}
}

func TestGetAllUser(t *testing.T) {
	users, err := GetAllUsers(db)
	if err != nil {
		t.Errorf("Error when getting all users. %v", err)
	}
	if len(users) == 0 {
		t.Errorf("No users found.")
	}
}

func TestAddUserNoFirstName(t *testing.T) {
	newUser := User{
		FirstName: "",
		LastName: templateUser.LastName,
		Gender: templateUser.Gender,
		Email: templateUser.Email,
		UserRole: templateUser.UserRole,
		Password: templateUser.Password,
	}
	err := testAddPatchRemoveUser(newUser)
	if err == nil {
		t.Errorf("User without first name was successfully added, patched and deleted")
	}
	matched, _ := regexp.MatchString("first_name", err.Error())
	if !matched {
		t.Errorf("Error did not contain any instance of first_name. %v", err)
	}
}

func TestAddUserNoLastName(t *testing.T) {
	newUser := User{
		FirstName: templateUser.FirstName,
		LastName: "",
		Gender: templateUser.Gender,
		Email: templateUser.Email,
		UserRole: templateUser.UserRole,
		Password: templateUser.Password,
	}
	err := testAddPatchRemoveUser(newUser)
	if err == nil {
		t.Errorf("User without last name was successfully added, patched and deleted")
	}
	matched, _ := regexp.MatchString("last_name", err.Error())
	if !matched {
		t.Errorf("Error did not contain any instance of last_name. %v", err)
	}
}

func TestAddUserNoGender(t *testing.T) {
	newUser := User{
		FirstName: templateUser.FirstName,
		LastName: templateUser.LastName,
		Gender: "",
		Email: templateUser.Email,
		UserRole: templateUser.UserRole,
		Password: templateUser.Password,
	}
	err := testAddPatchRemoveUser(newUser)
	if err == nil {
		t.Errorf("User without last name was successfully added, patched and deleted")
	}
	matched, _ := regexp.MatchString("gender", err.Error())
	if !matched {
		t.Errorf("Error did not contain any instance of gender. %v", err)
	}
}

func TestAddUserNoEmail(t *testing.T) {
	newUser := User{
		FirstName: templateUser.FirstName,
		LastName: templateUser.LastName,
		Gender: templateUser.Gender,
		Email: "",
		UserRole: templateUser.UserRole,
		Password: templateUser.Password,
	}
	err := testAddPatchRemoveUser(newUser)
	if err == nil {
		t.Errorf("User without email was successfully added, patched and deleted")
	}
	matched, _ := regexp.MatchString("email", err.Error())
	if !matched {
		t.Errorf("Error did not contain any instance of email. %v", err)
	}
}

func TestAddUserNoUserRole(t *testing.T) {
	newUser := User{
		FirstName: templateUser.FirstName,
		LastName: templateUser.LastName,
		Gender: templateUser.Gender,
		Email: templateUser.Email,
		UserRole: "",
		Password: templateUser.Password,
	}
	err := testAddPatchRemoveUser(newUser)
	if err == nil {
		t.Errorf("User without user_role was successfully added, patched and deleted")
	}
	matched, _ := regexp.MatchString("user_role", err.Error())
	if !matched {
		t.Errorf("Error did not contain any instance of user_role. %v", err)
	}
}

func TestAddUserValid(t *testing.T) {
	if err := testAddPatchRemoveUser(templateUser); err != nil {
		t.Errorf("Something went wrong with a valid user. %v", err)
	}
}