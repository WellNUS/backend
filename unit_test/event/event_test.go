package event

import (
	"wellnus/backend/db/model"
	"wellnus/backend/unit_test/test_helper"

	"testing"
	"net/http"
	"fmt"
)

// Full test
func TestEventHandler(t *testing.T) {
	t.Run("AddEventHandler no eventname", testAddEventHandlerNoEventName)
	t.Run("AddEventHandler no access", testAddEventHandlerNoAccess)
	t.Run("AddEventHandler no category", testAddEventHandlerNoCategory)
	t.Run("AddEventHandler not logged in", testAddEventHandlerNotLoggedIn)
	t.Run("AddEventHandler as user1", testAddEventHandlerAsUser1)
	t.Run("GetAllEventsHandler as user0", testGetAllEventsHandlerAsUser0)
	t.Run("GetAllEventsHandler not logged in", testGetAllEventsHandlerAsNotLoggedIn)
	t.Run("GetAllEventsHandler as user1", testGetAllEventsHandlerAsUser1)
	t.Run("GetEvent0Handler not logged in", testGetEvent0HandlerAsNotLoggedIn)
	t.Run("GetEvent0Handler as user1", testGetEvent0HandlerAsUser1)
	t.Run("AddUser1ToEvent0 as user1", testAddUser1ToPublicEvent0HandlerAsUser1)
	t.Run("AddUser0ToEvent1 as user0", testAddUser0ToPrivateEvent1HandlerAsUser0)
	t.Run("AddUser0ToEvent1 not logged in", testAddUser0ToPrivateEvent1HandlerNotLoggedIn)
	t.Run("Adduser0ToEvent1 as user1", testAddUser0ToPrivateEvent1HandlerAsUser1)
	t.Run("GetAllEventHandler as user0 after addition", testGetAllEventHandlerAsUser0AfterAddition)
	t.Run("UpdateEvent0Handler as not user0", testUpdateEvent0HandlerAsNotUser0)
	t.Run("UpdateEvent0Handler as user0", testUpdateEvent0HandlerAsUser0)
	t.Run("GetAllEventshandler as user1 after update", testGetAllEventsHandlerAsUser1AfterUpdate)
	t.Run("LeaveEvent1Handler as user0", testLeaveEvent1HandlerAsUser0)
	t.Run("LeaveEvent1handler as user1", testLeaveEvent1HandlerAsUser1)
	t.Run("LeaveAllEventsHandler as user0", testLeaveAllEventsHandlerAsUser0)
	t.Run("CreateGroupDeleteEventHandler Unauthorised", testCreateGroupDeleteEventUnauthorised)
	t.Run("CreateGroupDeleteEventHandler Authorised", testCreateGroupDeleteEventAuthorised)
}

func testAddEventHandlerNoEventName(t *testing.T) {
	testEvent := test_helper.GetTestEvent(1)
	testEvent.EventName = ""
	ioReaderEvent, _ := test_helper.GetIOReaderFromObject(testEvent)
	req, _ := http.NewRequest("POST", "/event", ioReaderEvent)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code == http.StatusOK {
		t.Errorf("Event with no event_name sucessfully added. Status Code: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, "event_name")
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of event_name. %s", errString)
	}
}

func testAddEventHandlerNoAccess(t *testing.T) {
	testEvent := test_helper.GetTestEvent(1)
	testEvent.Access = ""
	ioReaderEvent, _ := test_helper.GetIOReaderFromObject(testEvent)
	req, _ := http.NewRequest("POST", "/event", ioReaderEvent)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code == http.StatusOK {
		t.Errorf("Event with no access sucessfully added. Status Code: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, "access")
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of access. %s", errString)
	}
}

func testAddEventHandlerNoCategory(t *testing.T) {
	testEvent := test_helper.GetTestEvent(1)
	testEvent.Category = ""
	ioReaderEvent, _ := test_helper.GetIOReaderFromObject(testEvent)
	req, _ := http.NewRequest("POST", "/event", ioReaderEvent)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code == http.StatusOK {
		t.Errorf("Event with no category sucessfully added. Status Code: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, "category")
	if !matched {
		t.Errorf("response body was not an error did not contain any instance of category. %s", errString)
	}
}

func testAddEventHandlerNotLoggedIn(t *testing.T) {
	ioReaderEvent, _ := test_helper.GetIOReaderFromObject(test_helper.GetTestEvent(1))
	req, _ := http.NewRequest("POST", "/event", ioReaderEvent)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code == http.StatusOK {
		t.Errorf("Event was sucessfully added when not logged in. Status Code: %d", w.Code)
	}
}

func testAddEventHandlerAsUser1(t *testing.T) {
	ioReaderEvent, _ := test_helper.GetIOReaderFromObject(test_helper.GetTestEvent(1))
	req, _ := http.NewRequest("POST", "/event", ioReaderEvent)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to AddEvent failed with status code of %d", w.Code)
	}
	eventWithUsers, err := test_helper.GetEventWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting testEvent1 from body. %v", err)
	}
	testEvents = append(testEvents, eventWithUsers.Event)
	if testEvents[1].ID == 0 {
		t.Errorf("testEvent1 ID was not written by addEvent call")
	}
	if testEvents[1].OwnerID != testUsers[1].ID {
		t.Errorf("testUser1 is not owner of event despite being the one who created event")
	}
	if len(eventWithUsers.Users) < 1 || eventWithUsers.Users[0].ID != testUsers[1].ID {
		t.Errorf("Owner was not added into the new event")
	}
	// fmt.Printf("%v \n", testEvents)
}

func testGetAllEventsHandlerAsUser0(t *testing.T) {
	req, _ := http.NewRequest("GET", "/event", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetAllEvent failed with status code of %d", w.Code)
	}
	user0Events, err := test_helper.GetEventsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all events of user0 from body. %v", err)
	}
	if l := len(user0Events); l != 1 {
		t.Errorf("GetAllEventsHandler does not show 1 event but instead shows %d events", len(user0Events))
	}
}

func testGetAllEventsHandlerAsNotLoggedIn(t *testing.T) {
	req, _ := http.NewRequest("GET", "/event", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetAllEvents failed with status code of %d", w.Code)
	}
	events, err := test_helper.GetEventsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all events of not logged in. %v", err)
	}
	if l := len(events); l != 0 {
		t.Errorf("GetAllEventsHandler does not show 0 events but instead shows %d events", len(events))
	}
}

func testGetAllEventsHandlerAsUser1(t *testing.T) {
	req, _ := http.NewRequest("GET", "/event", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetAllEvents failed with status code of %d", w.Code)
	}
	user1Events, err := test_helper.GetEventsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all events of user1 . %v", err)
	}
	if l := len(user1Events); l != 1 {
		t.Errorf("GetAllEventsHandler does not show 1 events but instead shows %d events", len(user1Events))
	}
}

func testGetEvent0HandlerAsNotLoggedIn(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/event/%d", testEvents[0].ID), nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetEvent failed with status code of %d", w.Code)
	}
	eventWithUsers, err := test_helper.GetEventWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting event with users of not logged in. %v", err)
	}
	if l := len(eventWithUsers.Users); l != 1 {
		t.Errorf("The number of users in event is %d and not 1", l)
	}
	if id := eventWithUsers.Users[0].ID; id != testUsers[0].ID {
		t.Errorf("The user in the event is not user0 but user with ID = %d", id)
	}
}

func testGetEvent0HandlerAsUser1(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/event/%d", testEvents[0].ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetEvent failed with status code of %d", w.Code)
	}
	eventWithUsers, err := test_helper.GetEventWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting event with users of not logged in. %v", err)
	}
	if l := len(eventWithUsers.Users); l != 1 {
		t.Errorf("The number of users in event is %d and not 1", l)
	}
	if id := eventWithUsers.Users[0].ID; id != testUsers[0].ID {
		t.Errorf("The user in the event is not user0 but user with ID = %d", id)
	}
}

func testAddUser1ToPublicEvent0HandlerAsUser1(t *testing.T) {
	ioReaderUserIDBody, _ := test_helper.GetIOReaderFromObject(UserIDBody{ UserID: testUsers[1].ID })
	req, _ := http.NewRequest("POST", fmt.Sprintf("/event/%d", testEvents[0].ID), ioReaderUserIDBody)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request failed with a status code: %d", w.Code)
	}
	eventWithUsers, err := test_helper.GetEventWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured when getting eventWithUsers. %v", err)
	}
	if l := len(eventWithUsers.Users); l != 2 {
		t.Errorf("User1 was not added to event0. event0 has %d users", l)
	}
}

func testAddUser0ToPrivateEvent1HandlerAsUser0(t *testing.T) {
	ioReaderUserIDBody, _ := test_helper.GetIOReaderFromObject(UserIDBody{ UserID: testUsers[1].ID })
	req, _ := http.NewRequest("POST", fmt.Sprintf("/event/%d", testEvents[1].ID), ioReaderUserIDBody)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request did not fail with unauthorized but with code: %d", w.Code)
	}
}

func testAddUser0ToPrivateEvent1HandlerNotLoggedIn(t *testing.T) {
	ioReaderUserIDBody, _ := test_helper.GetIOReaderFromObject(UserIDBody{ UserID: testUsers[1].ID })
	req, _ := http.NewRequest("POST", fmt.Sprintf("/event/%d", testEvents[1].ID), ioReaderUserIDBody)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request did not fail with unauthorized but with code: %d", w.Code)
	}
}

func testAddUser0ToPrivateEvent1HandlerAsUser1(t *testing.T) {
	ioReaderUserIDBody, _ := test_helper.GetIOReaderFromObject(UserIDBody{ UserID: testUsers[0].ID })
	req, _ := http.NewRequest("POST", fmt.Sprintf("/event/%d", testEvents[1].ID), ioReaderUserIDBody)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request failed with a status code: %d", w.Code)
	}
	eventWithUsers, err := test_helper.GetEventWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured when getting eventWithUsers. %v", err)
	}
	if l := len(eventWithUsers.Users); l != 2 {
		t.Errorf("User0 was not added to event1. event1 has %d users", l)
	}
}

func testGetAllEventHandlerAsUser0AfterAddition(t *testing.T) {
	req, _ := http.NewRequest("GET", "/event", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetAllEvents failed with status code of %d", w.Code)
	}
	user0Events, err := test_helper.GetEventsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting all events of user0. %v", err)
	}
	if l := len(user0Events); l != 2 {
		t.Errorf("GetAllEventsHandler does not show 2 events but instead shows %d events", l)
	}
}

func testUpdateEvent0HandlerAsNotUser0(t *testing.T) {
	ioReaderEvent, _ := test_helper.GetIOReaderFromObject(Event{ EventName: "UpdatedEventName" })
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/event/%d", testEvents[0].ID), ioReaderEvent)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request did not give an Unauthorized status code. Status Code: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, UnauthorizedErrorMessage)
	if !matched {
		t.Errorf("User that was not logged in was not unauthorised. %s", errString)
	}

	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w = test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request did not give an Unauthorized status code. Status Code: %d", w.Code)
	}
	errString, matched = test_helper.CheckErrorMessageFromRecorder(w, UnauthorizedErrorMessage)
	if !matched {
		t.Errorf("Unauthorized user was not unauthorized. %s", errString)
	}
}

func testUpdateEvent0HandlerAsUser0(t *testing.T) {
	newEventName := "UpdatedEventName"
	ioReaderEvent, _ := test_helper.GetIOReaderFromObject(Event{ EventName: newEventName })
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/event/%d", testEvents[0].ID), ioReaderEvent)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request did not give an OK status code. Status Code: %d", w.Code)
	}
	event, err := test_helper.GetEventFromRecorder(w)
	if err != nil {
		t.Errorf("There was an error while getting the event from body. %v", err)
	}
	if event.EventName != newEventName {
		t.Errorf("Returned event did not have the updated event_name")
	}
}

func testGetAllEventsHandlerAsUser1AfterUpdate(t *testing.T) {
	req, _ := http.NewRequest("GET", "/event", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to GetEvent failed with status code of %d", w.Code)
	}
	user1Events, err := test_helper.GetEventsFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting user1 events from body. %v", err)
	}
	if l := len(user1Events); l != 2 {
		t.Errorf("The number of user 2 events is %d and not 2", l)
	}
	for _, event := range user1Events {
		if event.ID == testEvents[0].ID && event.EventName != "UpdatedEventName" {
			t.Errorf("The event name was not updated from previous test and is instead %s", event.EventName)
		}
	}
	
}

func testLeaveEvent1HandlerAsUser0(t *testing.T) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/event/%d", testEvents[1].ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to LeaveEvent failed with status code of %d", w.Code)
	}
	eventWithUsers, err := test_helper.GetEventWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while getting event with users from body. %v", err)
	}
	if users := eventWithUsers.Users; len(users) != 1 {
		t.Errorf("There was not 1 user remaining in the event. number of users in the event = %d", len(users))
	}
	if lastUser := eventWithUsers.Users[0]; lastUser.ID != testUsers[1].ID {
		t.Errorf("The last user in the event is not user1")
	}
}

func testLeaveEvent1HandlerAsUser1(t *testing.T) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/event/%d", testEvents[1].ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to LeaveEvent failed with status code of %d", w.Code)
	}

	rows, err := DB.Query("SELECT COUNT(*) FROM wn_event WHERE id = $1", testEvents[1].ID)
	if err != nil {
		t.Errorf("An error occured while querying count. %v", err)
	}
	count, err := test_helper.ReadInt(rows)
	if err != nil {
		t.Errorf("An error occured while retrieving count. %v", err)
	}
	if count != 0 {
		t.Errorf("Event1 was not deleted as count is %d", count)
	}
}

func testLeaveAllEventsHandlerAsUser0(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/event", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to LeaveAllEvents failed with status code of %d", w.Code)
	}
	rows, err := DB.Query("SELECT COUNT(*) FROM wn_event")
	if err != nil {
		t.Errorf("An error occured while querying count. %v", err)
	}
	count, err := test_helper.ReadInt(rows)
	if err != nil {
		t.Errorf("An error occured while retrieving count. %v", err)
	}
	if count != 0 {
		t.Errorf("Not all events were deleted was not deleted as count is %d", count)
	}
}

func testCreateGroupDeleteEventUnauthorised(t *testing.T) {
	testEvents, _ = test_helper.SetupEventForUsers(DB, testUsers[:1])
	eventWithUsers, _ := model.AddUserToEventAuthorized(DB, testUsers[1].ID, testEvents[0].ID, testUsers[0].ID)
	if l := len(eventWithUsers.Users); l != 2 {
		t.Errorf("Something went wrong wile setting up event for creategroupdeleteevent. group has only %d users", l)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("/event/%d/start", testEvents[0].ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to CreateGroupDeleteEvent did not give unauthorized status. StatusCode: %d", w.Code)
	}
	errString, matched := test_helper.CheckErrorMessageFromRecorder(w, UnauthorizedErrorMessage)
	if !matched {
		t.Errorf("Error thrown did not contain instance of 401 Unauthorized. %s", errString)
	}
}

func testCreateGroupDeleteEventAuthorised(t *testing.T) {
	testEvents, _ = test_helper.SetupEventForUsers(DB, testUsers[:1])
	eventWithUsers, _ := model.AddUserToEventAuthorized(DB, testUsers[1].ID, testEvents[0].ID, testUsers[0].ID)
	if l := len(eventWithUsers.Users); l != 2 {
		t.Errorf("Something went wrong wile setting up event for creategroupdeleteevent. group has only %d users", l)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("/event/%d/start", testEvents[0].ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to CreateGroupDeleteEvent did not give unauthorized status. StatusCode: %d", w.Code)
	}
	groupWithUsers, err := test_helper.GetGroupWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving groupWithUsers. %v", err)
	}
	if groupWithUsers.Group.OwnerID != testUsers[0].ID {
		t.Errorf("new group owner is not the event owner")
	}
	for i, refUser := range testUsers {
		found := false
		for _, user := range groupWithUsers.Users {
			if user.ID == refUser.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("user%d was not added to the group", i)
		}
	}

	_, err = model.GetEvent(DB, testEvents[0].ID)
	if err.Error() != NotFoundErrorMessage {
		t.Errorf("error thrown is same as NotfoundError message, suggesting event was not deleted. %v", err)
	}
}