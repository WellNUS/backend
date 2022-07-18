package counsel

import (
	"wellnus/backend/unit_test/test_helper"

	"testing"
	"net/http"
	"fmt"
	"time"
)

func TestCounselHandler(t *testing.T) {
	t.Run("Asserting user roles", testUserRoles)
	t.Run("Asserting counsel requests topics", testCounselRequestTopics)
	t.Run("AddCounselRequest as not logged in", testAddCounselRequestHandlerAsNotLoggedIn)
	t.Run("AddCounselRequest as User 2", testAddCounselRequestHandlerAsUser2)
	t.Run("GetAllCounselRequests as not logged in", testGetAllCounselRequestsHandlerAsNotLoggedIn)
	t.Run("GetAllCounselRequests as member", testGetAllCounselRequestsHandlerAsMember)
	t.Run("GetAllCounselRequests as volunteer", testGetAllCounselRequestHandlerAsVolunteer)
	t.Run("GetAllCounselRequests as counsellor", testGetAllCounselRequestHandlerAsCounsellor)
	t.Run("GetCounselRequests of Anxiety", testGetCounselRequestsHandlerOfAnxiety)
	t.Run("GetCounselRequests of Anxiety and offMyChest", testGetCounselRequestsHandlerOfAnxietyAndOffMyChest)
	t.Run("GetCounselRequests of Anxiety, OffMyChest and SelfHarm", testGetCounselRequestsHandlerOfAnxietyOffMyChestAndSelfHarm)
	t.Run("GetCounselRequest of User 0 as User 2", testGetCounselRequestHandlerOfUser0AsUser2)
	t.Run("UpdateCounselRequest of User0 No topics", testUpdateCounselRequestHandlerOfUser0NoTopics)
	t.Run("UpdateCounselRequest of User0 No time", testUpdateCounselRequestHandlerOfUser0NoTimeAdded)
	t.Run("UpdateCounselRequest of User0 Success", testUpdateCounselRequestHandlerOfUser0Successful)
	t.Run("DeleteCounselRequest as not logged in", testDeleteCounselRequestHandlerAsNotLoggedIn)
	t.Run("DeleteCounselRequest as user1", testDeleteCounselRequestHandlerAsUser1)
	t.Run("AcceptCounselRequest of User2 as User0", testAcceptCounselRequestOfUser2AsUser0)
	t.Run("AcceptCounselRequest of user0 as User2", testAcceptCounselRequestOfUser0AsUser2)
	t.Run("GetCounselRequest of User0 after accept", testGetCounselRequestOfUser0AfterAccept)
}

func testUserRoles(t *testing.T) {
	if testUsers[0].UserRole != "MEMBER" {
		t.Errorf("User0 is not a member")
	}
	if testUsers[1].UserRole != "VOLUNTEER" {
		t.Errorf("User1 is not a volunteer")
	}
	if testUsers[2].UserRole != "COUNSELLOR" {
		t.Errorf("User2 is not a counsellor")
	}
}

func testCounselRequestTopics(t *testing.T) {
	if testCounselRequests[0].Topics[0] != "Anxiety" || testCounselRequests[0].Topics[1] != "OffMyChest" {
		t.Errorf("User0 counsel request topics did not match the test requirements")
	}
	if testCounselRequests[1].Topics[0] != "OffMyChest" || testCounselRequests[1].Topics[1] != "SelfHarm" {
		t.Errorf("User1 counsel request topics did not match the test requirements")
	}
}

func testAddCounselRequestHandlerAsNotLoggedIn(t *testing.T) {
	ioReaderCounselRequest, _ := test_helper.GetIOReaderFromObject(test_helper.GetTestCounselRequest(2))
	req, _ := http.NewRequest("POST", "/counsel", ioReaderCounselRequest)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to AddUpdateCounselRequest as not logged in did not have status code: %d", w.Code)
	}
}

func testAddCounselRequestHandlerAsUser2(t *testing.T) {
	ioReaderCounselRequest, _ := test_helper.GetIOReaderFromObject(test_helper.GetTestCounselRequest(2))
	req, _ := http.NewRequest("POST", "/counsel", ioReaderCounselRequest)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[2],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to AddUpdateCounselRequest as user 2 failed with status code: %d", w.Code)
	}
	counselRequest, err := test_helper.GetCounselRequestFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while add updating counsel request. %v", err)
	}
	if counselRequest.UserID != testUsers[2].ID {
		t.Errorf("Returned counsel request is not attached to user2 but to user %d", counselRequest.UserID)
	}
	testCounselRequests = append(testCounselRequests, counselRequest)
}

func testGetAllCounselRequestsHandlerAsNotLoggedIn(t *testing.T) {
	req, _ := http.NewRequest("GET", "/counsel", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Reqeust to GetAllCounselRequest as not logged in did not give unauthorised code. Status Code: %d", w.Code)
	}
}

func testGetAllCounselRequestsHandlerAsMember(t *testing.T) {
	req, _ := http.NewRequest("GET", "/counsel", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Reqeust to GetAllCounselRequest as member did not give unauthorised code. Status Code: %d", w.Code)
	}
}

func testGetAllCounselRequestHandlerAsVolunteer(t *testing.T) {
	req, _ := http.NewRequest("GET", "/counsel", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetAllCounselRequest as Volunteer failed with Status Code: %d", w.Code)
	}
	counselRequests, err := test_helper.GetCounselRequestsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all counsel requests. %v", err)
	}
	if l := len(counselRequests); l != 3 {
		t.Errorf("GetAllCounselRequestHandler does not show 2 counsel request but %d", l)
	}
}

func testGetAllCounselRequestHandlerAsCounsellor(t *testing.T) {
	req, _ := http.NewRequest("GET", "/counsel", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[2],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Reqeust to GetAllCounselRequest as Volunteer failed with Status Code: %d", w.Code)
	}
	counselRequests, err := test_helper.GetCounselRequestsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all counsel requests. %v", err)
	}
	if l := len(counselRequests); l != 3 {
		t.Errorf("GetAllCounselRequestHandler does not show 2 counsel request but %d", l)
	}
}

func testGetCounselRequestsHandlerOfAnxiety(t *testing.T) {
	req, _ := http.NewRequest("GET", "/counsel?topic=Anxiety", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[2],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetAllCounselRequest with anxiety failed with Status Code: %d", w.Code)
	}
	counselRequests, err := test_helper.GetCounselRequestsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all counsel requests with anxiety topic. %v", err)
	}
	if l := len(counselRequests); l != 2 {
		t.Errorf("GetAllCounselRequestsHandler does not show 2 counsel request but %d", l)
	}
	for i, cr := range counselRequests {
		if l := len(cr.Topics); l != 2 {
			t.Errorf("Retrieved %dth counselRequest did not have 2 topics but %d topics", i, l)
		}
		if !cr.HasTopic("Anxiety") {
			t.Errorf("Retrieved %dth counselRequest did not have anxiety topic", i)
		}
	}
}

func testGetCounselRequestsHandlerOfAnxietyAndOffMyChest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/counsel?topic=Anxiety&topic=OffMyChest", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[2],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetCounselRequests with Anxiety and OffMyChest failed with Status Code: %d", w.Code)
	}
	counselRequests, err := test_helper.GetCounselRequestsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all counsel requests with Anxiety and OffMyChest. %v", err)
	}
	if l := len(counselRequests); l != 1 {
		t.Errorf("GetCounselRequestsHandler does not show 1 counsel request but %d", l)
	}
	cr := counselRequests[0]
	if l := len(cr.Topics); l != 2 {
		t.Errorf("Retrieved counselRequest did not have 2 topics but %d topics", l)
	}
	if !cr.HasTopic("Anxiety") || !cr.HasTopic("OffMyChest") {
		t.Errorf("Retrieved counselRequest did not have anxiety and offMyChest topic")
	}
}

func testGetCounselRequestsHandlerOfAnxietyOffMyChestAndSelfHarm(t *testing.T) {
	req, _ := http.NewRequest("GET", "/counsel?topic=Anxiety&topic=OffMyChest&topic=SelfHarm", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[2],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetCounselRequests with Anxiety, OffMyChest and SelfHarm failed with Status Code: %d", w.Code)
	}
	counselRequests, err := test_helper.GetCounselRequestsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all counsel requests with Anxiety, OffMyChest and SelfHarm. %v", err)
	}
	if l := len(counselRequests); l != 0 {
		t.Errorf("GetCounselRequestsHandler does not show 0 counsel request but %d", l)
	}
}

func testGetCounselRequestHandlerOfUser0AsUser2(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/counsel/%d", testUsers[0].ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[2],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetCounselRequest of user0 by user2 failed with status code: %d", w.Code)
	}
	counselRequest, err := test_helper.GetCounselRequestFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting counsel request. %v", err)
	}
	if counselRequest.UserID != testUsers[0].ID {
		t.Errorf("GetCounselRequestHandler does not belong to user0 but to userID = %d", counselRequest.UserID)
	}
}

func testUpdateCounselRequestHandlerOfUser0NoTopics(t *testing.T) {
	ioReaderCounselRequest, _ := test_helper.GetIOReaderFromObject(CounselRequest{
		Details: "This is an updated counselRequest",
		TimeAdded: time.Now(),
	})
	req, _ := http.NewRequest("POST", "/counsel", ioReaderCounselRequest)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code == http.StatusOK {
		t.Errorf("HTTP Request to UpdateCounselRequest of user0 were successful")
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, "topics")
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of topics. %s", errString)
	}
}

func testUpdateCounselRequestHandlerOfUser0NoTimeAdded(t *testing.T) {
	ioReaderCounselRequest, _ := test_helper.GetIOReaderFromObject(CounselRequest{
		Details: "This is an updated counselRequest",
		Topics: []string{"Depression"},
	})
	req, _ := http.NewRequest("POST", "/counsel", ioReaderCounselRequest)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code == http.StatusOK {
		t.Errorf("HTTP Request to UpdateCounselRequest of user0 were successful")
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, "time")
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of time. %s", errString)
	}
}

func testUpdateCounselRequestHandlerOfUser0Successful(t *testing.T) {
	ioReaderCounselRequest, _ := test_helper.GetIOReaderFromObject(CounselRequest{
		Details: "This is an updated counselRequest",
		Topics: []string{"Depression"},
		TimeAdded: time.Now(),
	})
	req, _ := http.NewRequest("POST", "/counsel", ioReaderCounselRequest)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to UpdateCounselRequest of user0 failed with status code: %d", w.Code)
	}
	counselRequest, err := test_helper.GetCounselRequestFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting counsel request. %v", err)
	}
	if counselRequest.UserID != testUsers[0].ID {
		t.Errorf("returned counsel request is not attached to user0 but to user %d", counselRequest.UserID)
	}
	testCounselRequests[0] = counselRequest
}

func testDeleteCounselRequestHandlerAsNotLoggedIn(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/counsel", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to DeleteCounselRequest of did not give an unauthorised status. Status Code: %d", w.Code)
	}
}

func testDeleteCounselRequestHandlerAsUser1(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/counsel", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to DeleteCounselRequest of failed witha Status Code: %d", w.Code)
	}
	counselRequest, err := test_helper.GetCounselRequestFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting counsel request. %v", err)
	}
	if counselRequest.UserID != testUsers[1].ID {
		t.Errorf("returned counsel request is not attached to user0 but to user %d", counselRequest.UserID)
	}
}

func testAcceptCounselRequestOfUser2AsUser0(t *testing.T) {
	req, _ := http.NewRequest("POST", fmt.Sprintf("/counsel/%d", testUsers[2].ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to AcceptCounselRequest did not give an unauthorized status code. Status Code: %d", w.Code)
	}
}

func testAcceptCounselRequestOfUser0AsUser2(t *testing.T) {
	req, _ := http.NewRequest("POST", fmt.Sprintf("/counsel/%d", testUsers[0].ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[2],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to AcceptCounselRequest failed with Status Code: %d", w.Code)
	}
	groupWithUsers, err := test_helper.GetGroupWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured when retrieving groupWithUsers")
	}
	if groupWithUsers.Group.Category != "COUNSEL" {
		t.Errorf("Created group is not a counsel group")
	}
	if l := len(groupWithUsers.Users); l != 2 {
		t.Errorf("New Group did not have 2 users but %d users", l)
	}
}

func testGetCounselRequestOfUser0AfterAccept(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/counsel/%d", testUsers[0].ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("HTTP Request to GetCounselRqeuest did not give not found status but status code: %d", w.Code)
	}
}