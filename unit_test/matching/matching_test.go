package matching

import (
	"wellnus/backend/config"
	"wellnus/backend/unit_test/test_helper"

	"testing"
	"net/http"
)

func TestMatching(t *testing.T) {
	t.Run("TestAssertInitialState of DB", testAssertInitialDatabaseState)
	t.Run("TestAddMatchRequestHandler as User0", testAddMatchRequestHandlerAsUser0)
	t.Run("TestAssertState of DB after User0", testAssertDatabaseStateAfterUser0)
	t.Run("TestAddMatchRequestHandler as User1", testAddMatchRequestHandlerAsUser1)
	t.Run("TestAssertState of DB after User1", testAssertDatabaseStateAfterUser1)
}

// Helper
func assertDatabaseState(t *testing.T, nUsers, nGroups, nMS, nMR int) {
	rows, _ := DB.Query(`SELECT COUNT(*) FROM wn_user`)
	if count, _ := test_helper.ReadInt(rows); count != nUsers {
		t.Errorf("The number of users is not equal to %d. No. of users = %d", nUsers, count)
	}
	rows, _ = DB.Query(`SELECT COUNT(*) FROM wn_group`)
	if count, _ := test_helper.ReadInt(rows); count != nGroups {
		t.Errorf("The number of groups initially present is not %d. No. of groups = %d", nGroups, count)
	}
	rows, _ = DB.Query(`SELECT COUNT(*) FROM wn_match_setting`)
	if count, _ := test_helper.ReadInt(rows); count != nMS {
		t.Errorf("The number of match settings is not %d. No. of match request = %d", nMS, count)
	}
	rows, _ = DB.Query(`SELECT COUNT(*) FROM wn_match_request`)
	if count, _ := test_helper.ReadInt(rows); count != nMR {
		t.Errorf("The number of match requests made is not %d. No. of match request = %d", nMR, count)
	}
}

func testAssertInitialDatabaseState(t *testing.T) {
	assertDatabaseState(t, config.MATCH_THRESHOLD, 0, config.MATCH_THRESHOLD, config.MATCH_THRESHOLD - 2)
}

func testAddMatchRequestHandlerAsUser0(t *testing.T) {
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
		t.Errorf("user_id of match request did not match user0")
	}
}

func testAssertDatabaseStateAfterUser0(t *testing.T) {
	assertDatabaseState(t, config.MATCH_THRESHOLD, 0, config.MATCH_THRESHOLD, config.MATCH_THRESHOLD - 1)
}

func testAddMatchRequestHandlerAsUser1(t *testing.T) {
	req, _ := http.NewRequest("POST", "/match", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("Add to match request failed with status code = %d", w.Code)
	}
	matchRequest, err := test_helper.GetMatchRequestFromRecorder(w)
	if err != nil {
		t.Errorf("There was an error while retrieving match request from buffer. %v", err)
	}
	if matchRequest.UserID != testUsers[1].ID {
		t.Errorf("user_id of match request did not match user0")
	}
}

func testAssertDatabaseStateAfterUser1(t *testing.T) {
	nGroups := config.MATCH_THRESHOLD / config.MATCH_GROUPSIZE
	nMR := config.MATCH_THRESHOLD % config.MATCH_GROUPSIZE
	assertDatabaseState(t, config.MATCH_THRESHOLD, nGroups, config.MATCH_THRESHOLD, nMR)
}