package join

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"database/sql"
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
	t.Run("GetJoinRequestHandler as not logged in", testGetLoadedJoinRequestHandlerAsNotLoggedIn)
	t.Run("GetAllJoinRequestHandler as not logged in", testGetAllJoinRequestHandlerAsNotLoggedIn)
	t.Run("GetAllJoinRequestHandler as user1", testGetAllJoinRequestHandlerAsUser1)
	t.Run("GetAllJoinRequestSentHandler as user1", testGetAllJoinRequestHandlerSentAsUser1)
	t.Run("GetAllJoinRequestReceivedHandler as user1", testGetAllJoinRequestHandlerReceivedAsUser1)
	t.Run("GetAllJoinRequestHandler as user 2", testGetAllJoinRequestHandlerAsUser2)
	t.Run("GetAllJoinRequestSentHandler as user 2", testGetAllJoinRequestHandlerSentAsUser2)
	t.Run("GetAllJoinRequestReceivedHandler as user2", testGetAllJoinRequestHandlerSentAsUser2)
	t.Run("RespondJoinRequestHandler reject not logged in", testRespondJoinRequestHandlerRejectNotLoggedIn)
	t.Run("RespondJoinRequestHandler reject as user1", testRespondJoinRequestHandlerRejectAsUser1)
	t.Run("RespondJoinRequestHandler approve as user1", testRespondJoinRequestHandlerApproveAsUser1)
	t.Run("DeleteJoinRequestHandler as user1", testDeleteJoinRequestHandlerAsUser1)
	t.Run("DeleteJoinRequestHandler as user2", testDeleteJoinRequestHandlerAsUser2)
	t.Run("GetJoinRequest after deletion", testGetLoadedJoinRequestHandlerAfterDeletion)
}

// Helper
func getBufferFromRecorder(w *httptest.ResponseRecorder) *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Result().Body)
	return buf
}

func getLoadedJoinRequestFromRecorder(w *httptest.ResponseRecorder) (LoadedJoinRequest, error) {
	buf := getBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return LoadedJoinRequest{}, errors.New(buf.String())
	}

	var loadedJoinRequest LoadedJoinRequest
	err := json.NewDecoder(buf).Decode(&loadedJoinRequest)
	if err != nil {
		return LoadedJoinRequest{}, err
	}
	return loadedJoinRequest, nil
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

func getJoinRequestRespondFromRecorder(w *httptest.ResponseRecorder) (JoinRequestRespond, error) {
	buf := getBufferFromRecorder(w)
	if w.Code != http.StatusOK {
		return JoinRequestRespond{}, errors.New(buf.String())
	}

	var joinRequestRespond JoinRequestRespond
	err := json.NewDecoder(buf).Decode(&joinRequestRespond)
	if err != nil {
		return JoinRequestRespond{}, err
	}
	return joinRequestRespond, nil
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

func readInt(row *sql.Rows) (int, error) {
	row.Next()
	var c int
	if err := row.Scan(&c); err != nil { return 0, err }
	return c, nil
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
}

func testGetLoadedJoinRequestHandlerAsNotLoggedIn(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/join/%d", addedJoinRequest.ID), nil)
	w := simulateRequest(req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetJoinRequest failed with status code of %d", w.Code)
	}
	retrievedLoadedJoinRequest, err := getLoadedJoinRequestFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving new join request from response. %v", err)
	}
	if !equal_joinRequest(retrievedLoadedJoinRequest.JoinRequest, addedJoinRequest) {
		t.Errorf("The retrieved JoinRequest component did not match the added join request")
	}
	if !equal_user(retrievedLoadedJoinRequest.User, validAddedUser2) {
		t.Errorf("The retrieved User component did not match the added join  user 2")
	}
	if !equal_group(retrievedLoadedJoinRequest.Group, validAddedGroup) {
		t.Errorf("The retrieved User component did not match the added join  user 2")
	}
}

func testGetAllJoinRequestHandlerAsNotLoggedIn(t *testing.T) {
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

func testGetAllJoinRequestHandlerAsUser1(t *testing.T) {
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

func testGetAllJoinRequestHandlerSentAsUser1(t *testing.T) {
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

func testGetAllJoinRequestHandlerReceivedAsUser1(t *testing.T) {
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

func testGetAllJoinRequestHandlerAsUser2(t *testing.T) {
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

func testGetAllJoinRequestHandlerSentAsUser2(t *testing.T) {
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

func testGetAllJoinRequestHandlerReceivedAsUser2(t *testing.T) {
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
	joinRequestRespond := JoinRequestRespond{ Approve: false }
	ioReaderRespond, _ := getIOReaderFromJoinRequestRespond(joinRequestRespond)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/join/%d", addedJoinRequest.ID), ioReaderRespond)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to respond while authorized gave Status code: %d", w.Code)
	}
	_, err := getJoinRequestRespondFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving join request from response, %v", err)
	}

	//Assert joinRequest deleted
	rows, err := db.Query(
		`SELECT COUNT(*) FROM wn_join_request 
		WHERE id = $1`,
		addedJoinRequest.ID)
	if err != nil {
		t.Errorf("An error occured while getting count from DB. %v", err)
	}
	c, err := readInt(rows)
	if err != nil {
		t.Errorf("An error occured while reading int. %v", err)
	}
	if c != 0 {
		t.Errorf("The join request still exist and has not been deleted")
	}


	//Assert user2 was not added to group
	rows, err = db.Query(
		`SELECT COUNT(*) FROM wn_user_group
		WHERE user_id = $1`,
		validAddedUser2.ID)
	if err != nil {
		t.Errorf("An errpr pccired while getting count from DB. %v", err)
	}
	c, err = readInt(rows)
	if err != nil {
		t.Errorf("An error occured while reading int. %v", err)
	}
	if c != 0 {
		t.Errorf("User2 is in some group despite being rejected")
	}
}

func testRespondJoinRequestHandlerApproveAsUser1(t *testing.T) {
	_, err := db.Exec(
		`INSERT INTO wn_join_request (
			id,
			user_id,
			group_id
		) values ($1, $2, $3)`,
		addedJoinRequest.ID,
		addedJoinRequest.UserID,
		addedJoinRequest.GroupID)
	if err != nil {
		t.Errorf("An error occured while brute adding the join request back. %v", err)
	}

	joinRequestRespond := JoinRequestRespond{ Approve: true }
	ioReaderRespond, _ := getIOReaderFromJoinRequestRespond(joinRequestRespond)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/join/%d", addedJoinRequest.ID), ioReaderRespond)
	req.AddCookie(&http.Cookie{
		Name: "id",
		Value: fmt.Sprintf("%d", validAddedUser1.ID),
	})
	w := simulateRequest(req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to respond while authorized gave Status code: %d", w.Code)
	}
	_, err = getJoinRequestRespondFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving join request from response, %v", err)
	}

	//Assert joinRequest deleted
	rows, err := db.Query(
		`SELECT COUNT(*) FROM wn_join_request 
		WHERE id = $1`,
		addedJoinRequest.ID)
	if err != nil {
		t.Errorf("An error occured while getting count from DB. %v", err)
	}
	c, err := readInt(rows)
	if err != nil {
		t.Errorf("An error occured while reading int. %v", err)
	}
	if c != 0 {
		t.Errorf("The join request still exist and has not been deleted")
	}

	//Assert user2 was not added to group
	rows, err = db.Query(
		`SELECT COUNT(*) FROM wn_user_group
		WHERE user_id = $1 AND group_id = $2`,
		validAddedUser2.ID,
		validAddedGroup.ID)
	if err != nil {
		t.Errorf("An error occured while getting count from DB. %v", err)
	}
	c, err = readInt(rows)
	if err != nil {
		t.Errorf("An error occured while reading int. %v", err)
	}
	if c != 1 {
		t.Errorf("User2 is not in the group despite being approved")
	}
}

func testDeleteJoinRequestHandlerAsUser1(t *testing.T) {
	_, err := db.Exec(
		`INSERT INTO wn_join_request (
			id,
			user_id,
			group_id
		) values ($1, $2, $3)`,
		addedJoinRequest.ID,
		addedJoinRequest.UserID,
		addedJoinRequest.GroupID)
	if err != nil {
		t.Errorf("An error occured while brute adding the join request back. %v", err)
	}

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

func testGetLoadedJoinRequestHandlerAfterDeletion(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/join/%d", addedJoinRequest.ID), nil)
	w := simulateRequest(req)
	if w.Code != http.StatusNotFound { 
		t.Errorf("HTTP Request to GetJoinRequest did not respond with NotFound Code but with status code of %d", w.Code)
	}
}