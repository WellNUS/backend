package provider

import (
	"wellnus/backend/unit_test/test_helper"

	"testing"
	"net/http"
	"fmt"
)

// Full test
func TestProviderHandler(t *testing.T) {
	t.Run("AddProviderSettingHandler not logged in", testAddProviderSettingHandlerNotLoggedIn)
	t.Run("AddProviderSettingHandler as Member", testAddProviderSettingHandlerAsMember)
	t.Run("AddProviderSettingHandler", testAddProviderSettingHandler)
	t.Run("GetAllProviderHandler", testGetAllProvidersHandler)
	t.Run("GetProvidersHandler Of OffMyChest", testGetProvidersHandlerOfOffMyChest)
	t.Run("GetProvidersHandler Of OffMyChest and SelfHarm", testGetProvidersHandlerOfAnxietyOffMyChestAndSelfHarm)
	t.Run("GetProvidersHandler of Anxiety and OffMyChest And SelfHarm", testGetProvidersHandlerOfAnxietyOffMyChestAndSelfHarm)
	t.Run("GetProviderWithEventsHandler not found", testGetProviderWithEventsHandlerNotFound)
	t.Run("GetProviderWithEventsHandler", testGetProviderWithEventsHandler)
	t.Run("UpdateProviderSettingHandler not logged in", testUpdateProviderSettingHandlerNotLoggedIn)
	t.Run("UpdateProviderSettingHandler", testUpdateProviderSettingHandler)
	t.Run("DeleteProviderSettingHandler unauthorised", testDeleteProviderSettingHandlerUnauthorized)
	t.Run("DeleteProviderSettingHandler successful", testDeleteProviderSettingHandlerSuccessful)
}

// Helper
func testAddProviderSettingHandlerNotLoggedIn(t *testing.T) {
	ioReaderProviderSetting, _ := test_helper.GetIOReaderFromObject(test_helper.GetTestProviderSetting(0))
	req, _ := http.NewRequest("POST", "/provider", ioReaderProviderSetting)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to add get provider setting request not Unauthorized. Status Code: %d", w.Code)
	}
}

func testAddProviderSettingHandlerAsMember(t *testing.T) {
	ioReaderProviderSetting, _ := test_helper.GetIOReaderFromObject(test_helper.GetTestProviderSetting(0))
	req, _ := http.NewRequest("POST", "/provider", ioReaderProviderSetting)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to add get provider setting request not Unauthorized. Status Code: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, UnauthorizedErrorMessage)
	if !matched {
		t.Errorf("error did not have an instance of 401 Unauthorized. %s", errString )
	}
}

func testAddProviderSettingHandler(t *testing.T) {
	ioReaderProviderSetting, _ := test_helper.GetIOReaderFromObject(test_helper.GetTestProviderSetting(1))
	req, _ := http.NewRequest("POST", "/provider", ioReaderProviderSetting)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("Valid add providerSetting gave a status code = %d", w.Code)
	}
	var err error
	testProviderSettings[1], err = test_helper.GetProviderSettingFromRecorder(w)
	if err != nil {
			t.Errorf("There was an error while retrieving providerSetting from buffer. %v", err)
		}
	if testProviderSettings[1].UserID != testUsers[1].ID {
		t.Errorf("user_id of providerSetting did not match valid added user")
	}
}

func testGetAllProvidersHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/provider", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetAllProviders failed with Status Code: %d", w.Code)
	}
	providers, err := test_helper.GetProvidersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all providers topic. %v", err)
	}
	if l := len(providers); l != 2 {
		t.Errorf("GetAllProvidersHandler does not show 2 provider but %d", l)
	}
}

func testGetProvidersHandlerOfOffMyChest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/provider?topic=OffMyChest", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetProviders with offMyChest failed with Status Code: %d", w.Code)
	}
	providers, err := test_helper.GetProvidersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all providers with offMyChest topic. %v", err)
	}
	if l := len(providers); l != 2 {
		t.Errorf("GetAllProvidersHandler does not show 2 provider but %d", l)
	}
	for _, p := range providers {
		if !p.Setting.HasTopic("OffMyChest") {
			t.Errorf("Retrieved provider did not have offMyChest topic")
		}
	}
}

func testGetProvidersHandlerOfOffMyChestAndSelfHarm(t *testing.T) {
	req, _ := http.NewRequest("GET", "/provider?topic=SelfHarm&topic=OffMyChest", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetProviders with SelfHarm and OffMyChest failed with Status Code: %d", w.Code)
	}
	providers, err := test_helper.GetProvidersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all providers with SelfHarm and OffMyChest. %v", err)
	}
	if l := len(providers); l != 1 {
		t.Errorf("GetProvidersHandler does not show 1 provider but %d", l)
	}
	if !providers[0].Setting.HasTopic("SelfHarm") || !providers[0].Setting.HasTopic("OffMyChest") {
		t.Errorf("Retrieved provider did not have SelfHarm and offMyChest topic")
	}
}

func testGetProvidersHandlerOfAnxietyOffMyChestAndSelfHarm(t *testing.T) {
	req, _ := http.NewRequest("GET", "/provider?topic=Anxiety&topic=OffMyChest&topic=SelfHarm", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetProviders with Anxiety, OffMyChest and SelfHarm failed with Status Code: %d", w.Code)
	}
	providers, err := test_helper.GetProvidersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all providers with Anxiety, OffMyChest and SelfHarm. %v", err)
	}
	if l := len(providers); l != 0 {
		t.Errorf("GetProvidersHandler does not show 0 providers but %d", l)
	}
}

func testGetProviderWithEventsHandlerNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/provider/%d",testUsers[0].ID), nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("HTTP Request to get providerWithEvents does not give not NotFound. Status Code: %d", w.Code)
	}
}

func testGetProviderWithEventsHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/provider/%d",testUsers[1].ID), nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to get providerWithEvents failed with Status Code: %d", w.Code)
	}
	providerWithEvents, err := test_helper.GetProviderWithEventsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting providerWithEvents. %v", err)
	}
	if l := len(providerWithEvents.Events); l != 1 {
		t.Errorf("providerWithEvents does not contain 1 event but %d events", l)
	}
}

func testUpdateProviderSettingHandlerNotLoggedIn(t *testing.T) {
	ioReaderProviderSetting, _ := test_helper.GetIOReaderFromObject(test_helper.GetTestProviderSetting(0))
	req, _ := http.NewRequest("POST", "/provider", ioReaderProviderSetting)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Update to providerSetting failed with status code = %d", w.Code)
	}
}

func testUpdateProviderSettingHandler(t *testing.T) {
	ioReaderProviderSetting, _ := test_helper.GetIOReaderFromObject(test_helper.GetTestProviderSetting(0))
	req, _ := http.NewRequest("POST", "/provider", ioReaderProviderSetting)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("Update to providerSetting failed with status code = %d", w.Code)
	}
	var err error
	testProviderSettings[1], err = test_helper.GetProviderSettingFromRecorder(w)
	if err != nil {
		t.Errorf("There was an error while retrieving providerSetting from buffer. %v", err)
	}
	if !testProviderSettings[1].HasTopic("Anxiety") || !testProviderSettings[1].HasTopic("OffMyChest") {
		t.Errorf("ProviderSettings did not update to Anxiety and OffMyChest")
	}
}

func testDeleteProviderSettingHandlerUnauthorized(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/provider", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to delete providerSetting without loggedin was not unauthorised. Status Code: %d", w.Code)
	}
}

func testDeleteProviderSettingHandlerSuccessful(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/provider", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to delete providerSetting failed with a Status Code = %d", w.Code)
	}
	providerSetting, err := test_helper.GetProviderSettingFromRecorder(w)
	if err != nil {
		t.Errorf("There was an error while retrieving providerSetting from buffer. %v", err)
	}
	if providerSetting.UserID != testUsers[1].ID {
		t.Errorf("user_id of deleted providerSetting did not match id of user1")
	}
}