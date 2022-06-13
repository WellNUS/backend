package group

import (
	"testing"
	"regexp"
	"net/http"
	"net/http/httptest"
	"bytes"
	"errors"
	"encoding/json"
	"io"
	"fmt"
)

// Full test
func TestGroupHandler(t *testing.T) {
	t.Run("AddGroupHandler with no group name", testAddGroupHandlerNoGroupName)
	t.Run("AddGroupHandler with no category", testAddGroupHandlerNoCategory)
	t.Run("AddGroupHandler not logged in", testAddGroupHandlerNotLoggedIn)
	t.Run("AddGroupHandler successful as User1", testAddGroupHandlerAsUser1)
	t.Run("AddGroupHandler successful as User2 no description", testAddGroupHandlerAsUser2NoDescription)
	t.Run("GetAllGroupsHandler as User1", testGetAllGroupsHandlerAsUser1)
	t.Run("GetAllGroupsHandler as not logged in", testGetAllGroupsHandlerAsNotLoggedIn)
	t.Run("GetAllGroupsHandler as User2", testGetAllGroupsHandlerAsUser2)
	t.Run("GetGrouphandler not logged in", testGetGroupHandlerAsNotLoggedIn)
	t.Run("GetAllGroupHandler as User2 after joining", testGetAllGroupHandlerAsUser2AfterJoining)
	t.Run("UpdateGroupHandler as not User1", testUpdateGroupHandlerAsNotUser1)
	t.Run("UpdateGroupHandler as User1", testUpdateGroupHandlerAsUser1)
	t.Run("GetAllGroupHandler as User2", testGetAllGroupHandlerAsUser2)
	t.Run("LeaveGroupHandler as User1", testLeaveGroupHandlerAsUser1)
	t.Run("LeaveGroupHandler as User2", testLeaveAllGroupsHandlerAsUser2)
	t.Run("GetGrouphandler after delete", testGetGroupHandlerAfterDelete)
}

// Helper

func getBufferFromRecorder(w *httptest.ResponseRecorder) *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Result().Body)
	return buf
}

func getGroupFromRecorder(w *httptest.ResponseRecorder) (Group, error) {
	buf := getBufferFromRecorder(w)
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

func getGroupsFromRecorder(w *httptest.ResponseRecorder) ([]Group, error) {
	buf := getBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return nil, errors.New(buf.String())
	}
	// fmt.Printf("Response Body: %v \n", buf)
	groups := make([]Group, 0)
	err := json.NewDecoder(buf).Decode(&groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func getGroupWithUsersFromRecorder(w *httptest.ResponseRecorder) (GroupWithUsers, error) {
	buf := getBufferFromRecorder(w)
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

func getGroupsWithUsersFromRecorder(w *httptest.ResponseRecorder) ([]GroupWithUsers, error) {
	buf := getBufferFromRecorder(w)
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

func simulateRequest(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func getIOReaderFromGroup(group Group) (io.Reader, error) {
	j, err := json.Marshal(group)
	if err != nil { return nil, err }
	return bytes.NewReader(j), nil
}

func testAddGroupHandlerNoGroupName(t *testing.T) {
	testGroup := Group{
		GroupDescription: validAddedGroup1.GroupDescription,
		Category: validAddedGroup1.Category,
	}
	ioReaderGroup, _ := getIOReaderFromGroup(testGroup)
	req, _ := http.NewRequest("POST", "/group", ioReaderGroup)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	w := simulateRequest(req)
	if w.Code == http.StatusOK {
		t.Errorf("Group with no group_name sucessfully added. Status Code: %d", w.Code)
	}
	matched, _ := regexp.MatchString("group_name", getBufferFromRecorder(w).String())
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of group_name")
	}
}

func testAddGroupHandlerNoCategory(t *testing.T) {
	testGroup := Group{
		GroupName: validAddedGroup1.GroupName,
		GroupDescription: validAddedGroup1.GroupDescription,
	}
	ioReaderGroup, _ := getIOReaderFromGroup(testGroup)
	req, _ := http.NewRequest("POST", "/group", ioReaderGroup)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	w := simulateRequest(req)
	if w.Code == http.StatusOK {
		t.Errorf("Group with no category sucessfully added. Status Code: %d", w.Code)
	}
	matched, _ := regexp.MatchString("category", getBufferFromRecorder(w).String())
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of category")
	}
}

func testAddGroupHandlerNotLoggedIn(t *testing.T) {
	ioReaderGroup, _ := getIOReaderFromGroup(validAddedGroup1)
	req, _ := http.NewRequest("POST", "/group", ioReaderGroup)
	w := simulateRequest(req)
	if w.Code == http.StatusOK {
		t.Errorf("Group with no category sucessfully added. Status Code: %d", w.Code)
	}
}

func testAddGroupHandlerAsUser1(t *testing.T) {
	ioReaderGroup, _ := getIOReaderFromGroup(validAddedGroup1)
	req, _ := http.NewRequest("POST", "/group", ioReaderGroup)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	// fmt.Println("id:", validAddedUser1.ID)
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to AddGroup failed with status code of %d", w.Code)
	}
	var err error
	groupWithUsers, err := getGroupWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting validAddedGroup1 from body. %v", err)
	}
	validAddedGroup1 = groupWithUsers.Group
	if validAddedGroup1.ID == 0 {
		t.Errorf("validAddedGroup1 ID was not written by addGroup call")
	}
	if validAddedGroup1.OwnerID != validAddedUser1.ID {
		t.Errorf("validAddedUser1 is not owner of group despite being the one who created group")
	}
	if len(groupWithUsers.Users) < 1 || groupWithUsers.Users[0].ID != validAddedUser1.ID {
		t.Errorf("Owner was not added into the new group")
	}
}

func testAddGroupHandlerAsUser2NoDescription(t *testing.T) {
	ioReaderGroup, _ := getIOReaderFromGroup(validAddedGroup2)
	req, _ := http.NewRequest("POST", "/group", ioReaderGroup)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser2.ID),
	})
	// fmt.Println("id:", validAddedUser1.ID)
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to AddGroup failed with status code of %d", w.Code)
	}
	var err error
	groupWithUsers, err := getGroupWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting validAddedGroup2 from body. %v", err)
	}
	validAddedGroup2 = groupWithUsers.Group
	if validAddedGroup2.ID == 0 {
		t.Errorf("validAddedGroup2 ID was not written by addGroup call")
	}
	if validAddedGroup2.OwnerID != validAddedUser2.ID {
		t.Errorf("validAddedUser1 is not owner of group despite being the one who created group")
	}
	if len(groupWithUsers.Users) < 1 || groupWithUsers.Users[0].ID != validAddedUser2.ID {
		t.Errorf("Owner was not added into the new group")
	}
}

func testGetAllGroupsHandlerAsUser1(t *testing.T) {
	req, _ := http.NewRequest("GET", "/group", nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetAllGroup failed with status code of %d", w.Code)
	}
	user1Groups, err := getGroupsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all groups of user 1 from body. %v", err)
	}
	if l := len(user1Groups); l != 1 {
		t.Errorf("GetAllGroupsHandler does not show 1 group but instead shows %d groups", len(user1Groups))
	}
}

func testGetAllGroupsHandlerAsNotLoggedIn(t *testing.T) {
	req, _ := http.NewRequest("GET", "/group", nil)
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetAllGroups failed with status code of %d", w.Code)
	}
	groups, err := getGroupsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all groups of not logged in. %v", err)
	}
	if l := len(groups); l != 0 {
		t.Errorf("GetAllGroupsHandler does not show 0 groups but instead shows %d groups", len(groups))
	}
}

func testGetAllGroupsHandlerAsUser2(t *testing.T) {
	req, _ := http.NewRequest("GET", "/group", nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser2.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetAllGroups failed with status code of %d", w.Code)
	}
	user2Groups, err := getGroupsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all groups of user 2 . %v", err)
	}
	if l := len(user2Groups); l != 1 {
		t.Errorf("GetAllGroupsHandler does not show 1 groups but instead shows %d groups", len(user2Groups))
	}
}

func testGetGroupHandlerAsNotLoggedIn(t *testing.T) {
	req, _ := http.NewRequest("GET",fmt.Sprintf("/group/%d", validAddedGroup1.ID), nil)
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetGroup failed with status code of %d", w.Code)
	}
	groupWithUsers, err := getGroupWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting group with users of not logged int. %v", err)
	}
	if l := len(groupWithUsers.Users); l != 1 {
		t.Errorf("The number of users in group is %d and not 1", l)
	}
	if id := groupWithUsers.Users[0].ID; id != validAddedUser1.ID {
		t.Errorf("The user in the group is not user 1 but user with ID = %d", id)
	}
}

func testGetAllGroupHandlerAsUser2AfterJoining(t *testing.T) {
	_, err := db.Query(
		`INSERT INTO wn_user_group (
			user_id, 
			group_id) 
		VALUES ($1, $2)`, 
		validAddedUser2.ID, 
		validAddedGroup1.ID)
	if err != nil {
		t.Errorf("An error occured while adding user2 into group. %v", err)
	}

	req, _ := http.NewRequest("GET", "/group", nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser2.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetAllGroups failed with status code of %d", w.Code)
	}
	user2Groups, err := getGroupsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all groups of user 2 . %v", err)
	}
	if l := len(user2Groups); l != 2 {
		t.Errorf("GetAllGroupsHandler does not show 2 groups but instead shows %d groups", len(user2Groups))
	}
}

func testUpdateGroupHandlerAsNotUser1(t *testing.T) {
	ioReaderGroup, _ := getIOReaderFromGroup(Group{ GroupName: "UpdatedGroupName" })
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/group/%d", validAddedGroup1.ID), ioReaderGroup)
	w := simulateRequest(req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request did not give an Unauthorized status code. Status Code: %d", w.Code)
	}
	_, err := getGroupFromRecorder(w)
	match, _ := regexp.MatchString(UnauthorizedErrorMessage, err.Error())
	if !match {
		t.Errorf("User that was not logged in was not unauthorised. %v", err)
	}

	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser2.ID),
	})
	w = simulateRequest(req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request did not give an Unauthorized status code. Status Code: %d", w.Code)
	}
	_, err = getGroupFromRecorder(w)
	match, _ = regexp.MatchString(UnauthorizedErrorMessage, err.Error())
	if !match {
		t.Errorf("Unauthorized user was not unauthorised. %v", err)
	}
}

func testUpdateGroupHandlerAsUser1(t *testing.T) {
	newGroupName := "UpdatedGroupName"
	ioReaderGroup, _ := getIOReaderFromGroup(Group{ GroupName: newGroupName })
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/group/%d", validAddedGroup1.ID), ioReaderGroup)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request did not give an OK status code. Status Code: %d", w.Code)
	}
	group, err := getGroupFromRecorder(w)
	if err != nil {
		t.Errorf("There was an error while getting the group from body. %v", err)
	}
	if group.GroupName != newGroupName {
		t.Errorf("Returned group did not have the updated group_name")
	}
}

func testGetAllGroupHandlerAsUser2(t *testing.T) {
	req, _ := http.NewRequest("GET", "/group", nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser2.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetGroup failed with status code of %d", w.Code)
	}
	user2Groups, err := getGroupsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting user 2 groups from body. %v", err)
	}
	if l := len(user2Groups); l != 2 {
		t.Errorf("The number of user 2 groups is %d and not 2", l)
	}
	for _, group := range user2Groups {
		if group.ID == validAddedGroup1.ID && group.GroupName != "UpdatedGroupName" {
			t.Errorf("The group name was not updated from previous test and is instead %s", group.GroupName)
		}
	}
	
}

func testLeaveGroupHandlerAsUser1(t *testing.T) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/group/%d", validAddedGroup1.ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to LeaveGroup failed with status code of %d", w.Code)
	}
	groupWithUsers, err := getGroupWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting group with users from body. %v", err)
	}
	if ownerID := groupWithUsers.Group.OwnerID; ownerID != validAddedUser2.ID {
		t.Errorf("Ownership of group was not transferred to user 2")
	}
	if users := groupWithUsers.Users; len(users) != 1 {
		t.Errorf("There was not 1 user remaining in the group. number of users in the group = %d", len(users))
	}
	if lastUser := groupWithUsers.Users[0]; lastUser.ID != validAddedUser2.ID {
		t.Errorf("The last user in the group is not user 2")
	}
}

func testLeaveAllGroupsHandlerAsUser2(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/group", nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser2.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to LeaveAllGroups failed with status code of %d", w.Code)
	}
	groupsWithUsers, err := getGroupsWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting group with users from body. %v", err)
	}
	if l := len(groupsWithUsers); l != 2 {
		t.Errorf("The number of groups left was not 2 but %d", l)
	}
	if users := groupsWithUsers[0].Users; len(users) != 0 {
		t.Errorf("There was still remaining users in the group 1. number of users in the group = %d", len(users))
	}
	if users := groupsWithUsers[1].Users; len(users) != 0 {
		t.Errorf("There was still remaining users in the group 2. number of users in the group = %d", len(users))
	}
}

func testGetGroupHandlerAfterDelete(t *testing.T) {
	req, _ := http.NewRequest("GET",fmt.Sprintf("/group/%d", validAddedGroup1.ID), nil)
	w := simulateRequest(req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Group1 was not successfully deleted from prev test as indicated by code of %d", w.Code)
	}

	req, _ = http.NewRequest("GET",fmt.Sprintf("/group/%d", validAddedGroup2.ID), nil)
	w = simulateRequest(req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Group2 was not successfully deleted from prev test as indicated by code of %d", w.Code)
	}
}