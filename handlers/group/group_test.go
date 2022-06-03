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
	t.Run("AddGroupHandler successful as User1", testAddGroupHandlerSuccess)
	t.Run("GetAllGroupsHandler as User1", testGetAllGroupsHandlerAsUser1)
	t.Run("GetAllGroupsHandler as not User1", testGetAllGroupsHandlerAsNotUser1)
	t.Run("GetGrouphandler not logged in", testGetGroupHandlerAsNotLoggedIn)
	t.Run("GetAllGroupHandler as User2 after joining", testGetAllGroupHandlerAsUser2AfterJoining)
	t.Run("UpdateGroupHandler as not User1", testUpdateGroupHandlerAsNotUser1)
	t.Run("UpdateGroupHandler as User1", testUpdateGroupHandlerAsUser1)
	t.Run("GetAllGroupHandler as User2", testGetAllGroupHandlerAsUser2)
	t.Run("LeaveGroupHandler as User1", testLeaveGroupHandlerAsUser1)
	t.Run("LeaveGroupHandler as User2", testLeaveGroupHandlerAsUser2)
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
		GroupDescription: validAddedGroup.GroupDescription,
		Category: validAddedGroup.Category,
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
		GroupName: validAddedGroup.GroupName,
		GroupDescription: validAddedGroup.GroupDescription,
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
	ioReaderGroup, _ := getIOReaderFromGroup(validAddedGroup)
	req, _ := http.NewRequest("POST", "/group", ioReaderGroup)
	w := simulateRequest(req)
	if w.Code == http.StatusOK {
		t.Errorf("Group with no category sucessfully added. Status Code: %d", w.Code)
	}
}

func testAddGroupHandlerSuccess(t *testing.T) {
	ioReaderGroup, _ := getIOReaderFromGroup(validAddedGroup)
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
		t.Errorf("An error occured while getting validAddedGroup from body. %v", err)
	}
	validAddedGroup = groupWithUsers.Group
	if validAddedGroup.ID == 0 {
		t.Errorf("validAddedGroup ID was not written by addGroup call")
	}
	if validAddedGroup.OwnerID != validAddedUser1.ID {
		t.Errorf("validAddedUser1 is not owner of group despite being the one who created group")
	}
	if len(groupWithUsers.Users) < 1 || groupWithUsers.Users[0].ID != validAddedUser1.ID {
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

func testGetAllGroupsHandlerAsNotUser1(t *testing.T) {
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
	
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser2.ID),
	})
	w = simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetAllGroups failed with status code of %d", w.Code)
	}
	groups, err = getGroupsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all groups of user 2 . %v", err)
	}
	if l := len(groups); l != 0 {
		t.Errorf("GetAllGroupsHandler does not show 0 groups but instead shows %d groups", len(groups))
	}
}

func testGetGroupHandlerAsNotLoggedIn(t *testing.T) {
	req, _ := http.NewRequest("GET",fmt.Sprintf("/group/%d", validAddedGroup.ID), nil)
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
	err := addUserToGroup(db, validAddedGroup.ID, validAddedUser2.ID)
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
	if l := len(user2Groups); l != 1 {
		t.Errorf("GetAllGroupsHandler does not show 1 group but instead shows %d groups", len(user2Groups))
	}
}

func testUpdateGroupHandlerAsNotUser1(t *testing.T) {
	ioReaderGroup, _ := getIOReaderFromGroup(Group{ GroupName: "UpdatedGroupName" })
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/group/%d", validAddedGroup.ID), ioReaderGroup)
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
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/group/%d", validAddedGroup.ID), ioReaderGroup)
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
	if l := len(user2Groups); l != 1 {
		t.Errorf("The number of user 2 groups is %d and not 1", l)
	}
	if groupName := user2Groups[0].GroupName; groupName != "UpdatedGroupName" {
		t.Errorf("The group name was not updated from previous test and is instead %s", groupName)
	}
}

func testLeaveGroupHandlerAsUser1(t *testing.T) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/group/%d", validAddedGroup.ID), nil)
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

func testLeaveGroupHandlerAsUser2(t *testing.T) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/group/%d", validAddedGroup.ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser2.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to LeaveGroup failed with status code of %d", w.Code)
	}
	groupWithUsers, err := getGroupWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting group with users from body. %v", err)
	}
	if groupID := groupWithUsers.Group.ID; groupID != validAddedGroup.ID {
		t.Errorf("Returned group did not have the original groupID")
	}
	if users := groupWithUsers.Users; len(users) != 0 {
		t.Errorf("There was still remaining users in the group. number of users in the group = %d", len(users))
	}
}

func testGetGroupHandlerAfterDelete(t *testing.T) {
	req, _ := http.NewRequest("GET",fmt.Sprintf("/group/%d", validAddedGroup.ID), nil)
	w := simulateRequest(req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Group was not successfully deleted from prev test as indicated by code of %d", w.Code)
	}
}