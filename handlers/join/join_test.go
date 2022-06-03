package join

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
	"errors"
	"encoding/json"
	"io"
	"fmt"
)

var addedJoinRequest JoinRequest

// Full test
func TestJoinHandler(t *testing.T) {
	t.Run("AddJoinRequestHandler", testAddJoinRequestHandler)
	t.Run("GetJoinRequestHandler as not logged in", testGetJoinRequestAsNotLoggedIn)
	t.Run("GetAllJoinRequestHandler as not logged in", testGetAllJoinRequestAsNotLoggedIn)
	t.Run("GetAllJoinRequestHandler as user1", testGetAllJoinRequestAsUser1)
	t.Run("GetAllJoinRequestSentHandler as user1", testGetAllJoinRequestSentAsUser1)
	t.Run("GetAllJoinRequestReceivedHandler as user1", testGetAllJoinRequestReceivedAsUser1)
	t.Run("GetAllJoinRequestHandler as user 2", testGetAllJoinRequestAsUser2)
	t.Run("GetAllJoinRequestSentHandler as user 2", testGetAllJoinRequestSentAsUser2)
	t.Run("GetAllJoinRequestReceivedHandler as user2", testGetAllJoinRequestSentAsUser2)
	t.Run("RespondJoinRequestHandler reject not logged in", testRespondJoinRequestHandlerRejectNotLoggedIn)
	t.Run("RespondJoinRequestHandler reject as user1", testRespondJoinRequestHandlerRejectAsUser1)
	t.Run("RespondJoinRequestHandler approve as user1", testRespondJoinRequestHandlerApproveAsUser1)
	t.Run("DeleteJoinRequestHandler as user1", testDeleteJoinRequestHandlerAsUser1)
	t.Run("DeleteJoinRequestHandler as user2", testDeleteJoinRequestHandlerAsUser2)
	t.Run("GetJoinRequest after deletion", testGetJoinRequestAfterDeletion)
}

// Helper

func match(joinRequest1 JoinRequest, joinRequest2 JoinRequest) bool {
	return joinRequest1.ID == joinRequest2.ID &&
			joinRequest1.GroupID == joinRequest2.GroupID &&
			joinRequest1.UserID == joinRequest2.UserID &&
			joinRequest1.RequestStatus == joinRequest2.RequestStatus
}

func getBufferFromRecorder(w *httptest.ResponseRecorder) *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Result().Body)
	return buf
}

func getJoinRequestFromRecorder(w *httptest.ResponseRecorder) (JoinRequest, error) {
	buf := getBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return JoinRequest{}, errors.New(buf.String())
	}

	var joinRequest JoinRequest
	err := json.NewDecoder(buf).Decode(&joinRequest)
	if err != nil {
		return JoinRequest{}, err
	}
	return joinRequest, nil
}

func getJoinRequestsFromRecorder(w *httptest.ResponseRecorder) ([]JoinRequest, error) {
	buf := getBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return nil, errors.New(buf.String())
	}

	var joinRequests []JoinRequest
	err := json.NewDecoder(buf).Decode(&joinRequests)
	if err != nil {
		return nil, err
	}
	return joinRequests, nil
}

func simulateRequest(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func getIOReaderFromJoinRequestRespond(respond JoinRequestRespond) (io.Reader, error) {
	j, err := json.Marshal(respond)
	if err != nil { return nil, err }
	return bytes.NewReader(j), nil
}

func getIOReaderFromJoinRequest(joinRequest JoinRequest) (io.Reader, error) {
	j, err := json.Marshal(joinRequest)
	if err != nil { return nil, err }
	return bytes.NewReader(j), nil
}

func testAddJoinRequestHandler(t *testing.T) {
	ioReaderJoinRequest, err := getIOReaderFromJoinRequest(JoinRequest{ GroupID: validAddedGroup.ID })
	req, _ := http.NewRequest("POST", "/join", ioReaderJoinRequest)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser2.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to AddJoinRequest failed with status code of %d", w.Code)
	}
	addedJoinRequest, err = getJoinRequestFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving new join request from response. %v", err)
	}
	if addedJoinRequest.GroupID != validAddedGroup.ID {
		t.Errorf("Returned addedJoinRequest did not update one of its GroupID correctly")
	}
	if addedJoinRequest.UserID != validAddedUser2.ID {
		t.Errorf("Returned addedJoinRequest did not update one of its UserID correctly")
	}
	if addedJoinRequest.RequestStatus != "PENDING" {
		t.Errorf("Returned addedJoinRequest did not update one of its RequestStatus correctly")
	}
}

func testGetJoinRequestAsNotLoggedIn(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/join/%d", addedJoinRequest.ID), nil)
	w := simulateRequest(req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetJoinRequest failed with status code of %d", w.Code)
	}
	retrievedJoinRequest, err := getJoinRequestFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving new join request from response. %v", err)
	}
	if !match(retrievedJoinRequest, addedJoinRequest) {
		t.Errorf("The retrieved join request did not match the added join request")
	}
}

func testGetAllJoinRequestAsNotLoggedIn(t *testing.T) {
	req, _ := http.NewRequest("GET", "/join", nil)
	w := simulateRequest(req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllJoinRequest failed with status code of %d", w.Code)
	}
	joinRequests, err := getJoinRequestsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all join request from response. %v", err)
	}
	if len(joinRequests) != 0 {
		t.Errorf("A user who is not logged in saw some join requests directed to him")
	}
}

func testGetAllJoinRequestAsUser1(t *testing.T) {
	req, _ := http.NewRequest("GET", "/join", nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllJoinRequest failed with status code of %d", w.Code)
	}
	joinRequests, err := getJoinRequestsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all join request from response. %v", err)
	}
	if len(joinRequests) != 1 {
		t.Errorf("User1 does not see 1 join requests relevant to him")
	}
	if joinRequests[0].ID != addedJoinRequest.ID {
		t.Errorf("The single join request was not the added join request")
	}
}

func testGetAllJoinRequestSentAsUser1(t *testing.T) {
	req, _ := http.NewRequest("GET", "/join?request=SENT", nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllJoinRequest failed with status code of %d", w.Code)
	}
	joinRequests, err := getJoinRequestsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all join request from response. %v", err)
	}
	if len(joinRequests) != 0 {
		t.Errorf("User1 saw non-zero join request sent by it")
	}
}

func testGetAllJoinRequestReceivedAsUser1(t *testing.T) {
	req, _ := http.NewRequest("GET", "/join?request=RECEIVED", nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllJoinRequest failed with status code of %d", w.Code)
	}
	joinRequests, err := getJoinRequestsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all join request from response. %v", err)
	}
	if len(joinRequests) != 1 {
		t.Errorf("User1 does not see 1 join requests received by him")
	}
	if joinRequests[0].ID != addedJoinRequest.ID {
		t.Errorf("The single join request was not the added join request")
	}
}

func testGetAllJoinRequestAsUser2(t *testing.T) {
	req, _ := http.NewRequest("GET", "/join", nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser2.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllJoinRequest failed with status code of %d", w.Code)
	}
	joinRequests, err := getJoinRequestsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all join request from response. %v", err)
	}
	if len(joinRequests) != 1 {
		t.Errorf("User2 does not see 1 join requests relevant to it")
	}
	if joinRequests[0].ID != addedJoinRequest.ID {
		t.Errorf("The single join request was not the added join request")
	}
}

func testGetAllJoinRequestSentAsUser2(t *testing.T) {
	req, _ := http.NewRequest("GET", "/join?request=SENT", nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser2.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllJoinRequest failed with status code of %d", w.Code)
	}
	joinRequests, err := getJoinRequestsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all join request from response. %v", err)
	}
	if len(joinRequests) != 1 {
		t.Errorf("User2 does not see 1 join requests sent by it")
	}
	if joinRequests[0].ID != addedJoinRequest.ID {
		t.Errorf("The single join request was not the added join request")
	}
}

func testGetAllJoinRequestReceivedAsUser2(t *testing.T) {
	req, _ := http.NewRequest("GET", "/join?request=RECEIVED", nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser2.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllJoinRequest failed with status code of %d", w.Code)
	}
	joinRequests, err := getJoinRequestsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all join request from response. %v", err)
	}
	if len(joinRequests) != 0 {
		t.Errorf("User2 saw non-zero join request sent by it")
	}
}

func testRespondJoinRequestHandlerRejectNotLoggedIn(t *testing.T) {
	respond := JoinRequestRespond{ Approve: false }
	ioReaderRespond, _ := getIOReaderFromJoinRequestRespond(respond)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/join/%d", addedJoinRequest.ID), ioReaderRespond)
	w := simulateRequest(req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to respond while not logged in did not give status unauthorized but Status code: %d", w.Code)
	}
}

func testRespondJoinRequestHandlerRejectAsUser1(t *testing.T) {
	respond := JoinRequestRespond{ Approve: false }
	ioReaderRespond, _ := getIOReaderFromJoinRequestRespond(respond)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/join/%d", addedJoinRequest.ID), ioReaderRespond)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to respond while authorized gave Status code: %d", w.Code)
	}
	joinRequest, err := getJoinRequestFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving join request from response, %v", err)
	}
	if joinRequest.RequestStatus != "REJECTED" {
		t.Errorf("Returned join request did not have a status of REJECTED but %s", joinRequest.RequestStatus)
	}

	//Assert user2 was not added to group
	query := fmt.Sprintf(
		`SELECT COUNT(*) FROM wn_user_group
		WHERE user_id = %d`,
		validAddedUser2.ID)
	rows, err := db.Query(query)
	if err != nil {
		t.Errorf("An error occured while querying database. %v", err)
	}
	rows.Next()
	var c int
	if err := rows.Scan(&c); err != nil {
		t.Errorf("An error occured while reading row. %v", err)
	}
	if c != 0 {
		t.Errorf("User2 is in some group despite being rejected")
	}
}

func testRespondJoinRequestHandlerApproveAsUser1(t *testing.T) {
	respond := JoinRequestRespond{ Approve: true }
	ioReaderRespond, _ := getIOReaderFromJoinRequestRespond(respond)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/join/%d", addedJoinRequest.ID), ioReaderRespond)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to respond while authorized gave Status code: %d", w.Code)
	}
	joinRequest, err := getJoinRequestFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving join request from response, %v", err)
	}
	if joinRequest.RequestStatus != "APPROVED" {
		t.Errorf("Returned join request did not have a status of REJECTED but %s", joinRequest.RequestStatus)
	}

	//Assert user2 was not added to group
	query := fmt.Sprintf(
		`SELECT COUNT(*) FROM wn_user_group
		WHERE user_id = %d AND group_id = %d`,
		validAddedUser2.ID,
		validAddedGroup.ID)
	rows, err := db.Query(query)
	if err != nil {
		t.Errorf("An error occured while querying database. %v", err)
	}
	rows.Next()
	var c int
	if err := rows.Scan(&c); err != nil {
		t.Errorf("An error occured while reading row. %v", err)
	}
	if c != 1 {
		t.Errorf("User2 is not in the group despite being approved")
	}
}

func testDeleteJoinRequestHandlerAsUser1(t *testing.T) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/join/%d", addedJoinRequest.ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request did not respond with unauthorized code but gave Status code: %d", w.Code)
	}
}

func testDeleteJoinRequestHandlerAsUser2(t *testing.T) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/join/%d", addedJoinRequest.ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser2.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request did not respond with OK code but gave Status code: %d", w.Code)
	}
	joinRequest, err := getJoinRequestFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting join request from response. %v", err)
	}
	if joinRequest.ID != addedJoinRequest.ID {
		t.Errorf("Returned joinRequest did not have the ID of the added join request")
	}
}

func testGetJoinRequestAfterDeletion(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/join/%d", addedJoinRequest.ID), nil)
	w := simulateRequest(req)
	if w.Code != http.StatusNotFound { 
		t.Errorf("HTTP Request to GetJoinRequest did not respond with NotFound Code but with status code of %d", w.Code)
	}
}