package match

import (
	"wellnus/backend/unit_test/test_helper"

	"testing"
	"net/http"
	"fmt"
)

// Full test
func TestMatchHandler(t *testing.T) {
	t.Run("GetMatchRequestHandler Unauthorized", testGetMatchRequestHandlerUnauthorized)
	t.Run("GetMatchRequestHandler NotFound", testGetMatchRequestHandlerNotFound)
	t.Run("AddMatchRequestHandler NoSetting", testAddMatchRequestHandlerNoSetting)
	t.Run("GetMatchSettingHandler Unauthorized", testGetMatchSettingHandlerUnauthorized)
	t.Run("GetMatchSettingHandler NotFound", testGetMatchSettingHandlerNotFound)
	t.Run("AddMatchSettingHandler Unauthorized", testAddMatchSettingHandlerUnauthorized)
	t.Run("AddMatchSettingHandler Success", testAddMatchSettingHandler)
	t.Run("UpdateMatchSettingHandler Success", testUpdateMatchSettingHandler)
	t.Run("AddMatchRequestHandler Success", testAddMatchRequestHandlerSuccessful)
	t.Run("GetMatchRequestCountHandler", testGetMatchRequestCountHandler)
	t.Run("GetMatchRequestHandler Success", testGetMatchRequestHandlerSuccessful)
	t.Run("DeleteMatchRequestHandler Unauthorized", testDeleteMatchRequestHandlerUnauthorized)
	t.Run("DeleteMatchRequestHandler Success", testDeleteMatchRequestHandlerSuccessful)
	t.Run("DeleteMatchSettingHandler Unauthorized", testDeleteMatchSettingHandlerUnauthorized)
	t.Run("DeleteMatchSettingHandler Success", testDeleteMatchSettingHandlerSuccessful)
}

// Helper

func testGetMatchRequestHandlerUnauthorized(t *testing.T) {
	route := fmt.Sprintf("/match/%d", testUsers[0].ID)
	req, _ := http.NewRequest("GET", route, nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to get match request not unauthorized. Status Code: %d", w.Code)
	}
}

func testGetMatchRequestHandlerNotFound(t *testing.T) {
	route := fmt.Sprintf("/match/%d", testUsers[0].ID)
	req, _ := http.NewRequest("GET", route, nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("HTTP Request to get match request not NotFound. Status Code: %d", w.Code)
	}
}

func testAddMatchRequestHandlerNoSetting(t *testing.T) {
	req, _ := http.NewRequest("POST", "/match", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Adding match request without setting passed with status code = %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, "foreign key constraint")
	if !matched {
		t.Errorf("The error that occured was not due to match setting. %s", errString)
	}
}

func testGetMatchSettingHandlerUnauthorized(t *testing.T) {
	req, _ := http.NewRequest("GET", "/setting", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to get match setting request not unauthorized. Status Code: %d", w.Code)
	}
}

func testGetMatchSettingHandlerNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/setting", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("HTTP Request to get match setting request not NotFound. Status Code: %d", w.Code)
	}
}

func testAddMatchSettingHandlerUnauthorized(t *testing.T) {
	ioReaderMatchSetting, _ := test_helper.GetIOReaderFromObject(validMatchSetting)
	req, _ := http.NewRequest("POST", "/setting", ioReaderMatchSetting)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to add get match setting request not Unauthorized. Status Code: %d", w.Code)
	}
}

func testAddMatchSettingHandler(t *testing.T) {
	ioReaderMatchSetting, _ := test_helper.GetIOReaderFromObject(validMatchSetting)
	req, _ := http.NewRequest("POST", "/setting", ioReaderMatchSetting)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("Valid add Match setting gave a status code = %d", w.Code)
	}
	var err error
	validMatchSetting, err = test_helper.GetMatchSettingFromRecorder(w)
	if err != nil {
		t.Errorf("There was an error while retrieving matchSetting from buffer. %v", err)
	}
	if validMatchSetting.UserID != testUsers[0].ID {
		t.Errorf("user_id of match setting did not match valid added user")
	}
}

func testUpdateMatchSettingHandler(t *testing.T) {
	validMatchSetting.MBTI = "ESFJ"
	ioReaderMatchSetting, _ := test_helper.GetIOReaderFromObject(validMatchSetting)
	req, _ := http.NewRequest("POST", "/setting", ioReaderMatchSetting)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("Update to match setting failed with status code = %d", w.Code)
	}
	var err error
	validMatchSetting, err = test_helper.GetMatchSettingFromRecorder(w)
	if err != nil {
		t.Errorf("There was an error while retrieving matchSetting from buffer. %v", err)
	}
	if validMatchSetting.MBTI != "ESFJ" {
		t.Errorf("MBTI field of setting was not updated.")
	}
	if validMatchSetting.UserID != testUsers[0].ID {
		t.Errorf("user_id of match setting did not match valid added user")
	}
}

func testAddMatchRequestHandlerSuccessful(t *testing.T) {
	req, _ := http.NewRequest("POST", "/match", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("Add to match request failed with status code = %d", w.Code)
	}
	matchRequest, err := test_helper.GetMatchRequestFromRecorder(w)
	if err != nil {
		t.Errorf("There was an error while retrieving match request from buffer. %v", err)
	}
	if matchRequest.UserID != testUsers[0].ID{
		t.Errorf("user_id of match request did not match valid added user")
	}
}

func testGetMatchRequestCountHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/match", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("Get to match request failed with status code = %d", w.Code)
	}
	count, err := test_helper.GetInt64FromRecorder(w)
	if err != nil {
		t.Errorf("There was an error while retrieving int64 from buffer. %v", err)
	}
	if count != 1 {
		t.Errorf("The number of match request is not 1 but %d", count)
	}
}

func testGetMatchRequestHandlerSuccessful(t *testing.T) {
	route := fmt.Sprintf("/match/%d", testUsers[0].ID)
	req, _ := http.NewRequest("GET", route, nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("Get to match request failed with status code = %d", w.Code)
	}
	_, err := test_helper.GetLoadedMatchRequestFromRecorder(w)
	if err != nil {
		t.Errorf("There was an error while retrieving loaded match request from buffer. %v", err)
	}
}

func testDeleteMatchRequestHandlerUnauthorized(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/match", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to delete match request without loggedin was not unauthorised. Status Code: %d", w.Code)
	}
}

func testDeleteMatchRequestHandlerSuccessful(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/match", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to delete match request failed with a Status Code = %d", w.Code)
	}
	matchRequest, err := test_helper.GetMatchRequestFromRecorder(w)
	if err != nil {
		t.Errorf("There was an error while retrieving match request from buffer. %v", err)
	}
	if matchRequest.UserID != testUsers[0].ID {
		t.Errorf("user_id of deleted match request did not match id of valid added user")
	}
}

func testDeleteMatchSettingHandlerUnauthorized(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/setting", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to delete match setting without loggedin was not unauthorised. Status Code: %d", w.Code)
	}
}

func testDeleteMatchSettingHandlerSuccessful(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/setting", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to delete match setting failed with a Status Code = %d", w.Code)
	}
	matchSetting, err := test_helper.GetMatchSettingFromRecorder(w)
	if err != nil {
		t.Errorf("There was an error while retrieving match setting from buffer. %v", err)
	}
	if matchSetting.UserID != testUsers[0].ID {
		t.Errorf("user_id of deleted match setting did not match id of valid added user")
	}
}