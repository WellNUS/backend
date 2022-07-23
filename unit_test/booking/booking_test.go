package booking

import (
	"wellnus/backend/unit_test/test_helper"
	"testing"
	"net/http"
	"fmt"
)

// Full test
func TestBookingHandler(t *testing.T) {
	t.Run("AddBookingHandler to no provider setting", testAddBookingHandlerToNoProviderSettingUser1ByUser0)
	t.Run("GetBookingProviderHandler not logged in of no provider setting", testGetBookingProviderHandlerAsNotLoggedInNoProviderSetting)
	t.Run("GetAllBookingUsersHandler not logged in", testGetAllBookingUsersHandlerAsNotLoggedIn)
	t.Run("GetAllBookingUsersHander as user1", testGetAllBookingUserHandlerAsUser1)
	t.Run("GetAllBookingUsersHandler sent as user1", testGetAllBookingUsersHandlerSentAsUser1)
	t.Run("GetAllBookingUsersHandler received as user1", testGetAllBookingUsersHandlerReceivedAsUser1)
	t.Run("GetAllBookingUsersHandler required as user1", testGetAllBookingUsersHandlerRequiredAsUser1)
	t.Run("GetAllBookingusershandler required as user2", testGetAllBookingUsersHandlerRequiredAsUser2)
	t.Run("RespondBookingHandler reject not logged in", testRespondBookingHandlerRejectNotLoggedIn)
	t.Run("RespondBookingHandler reject as user1", testRespondBookingHandlerRejectAsUser1)
	t.Run("GetAllBookingUsersHandler required as user0", testGetAllBookingUserHandlerRequiredAsUser0)
	t.Run("RespondBookingHandler approve as user0", testRespondBookingHandlerApproveAsUser0)
	t.Run("UpdateBookinghandler not logged in", testUpdateBookingHandlerOfBooking1To2NotLoggedIn)
	t.Run("UpdateBookingHandler unauthorised", testUpdateBookingHandlerOfBooking1To2AsUser2Unauthorized)
	t.Run("UpdateBookingHandler authorised", testUpdateBookingHandlerOfBooking1To2AsUser1Authorized)
	t.Run("GetBookingProvider after update", testGetBookingProviderOfBooking1To2AfterUpdate)
	t.Run("DeleteBookingHandler unauthorised", testDeleteBookingHandlerOfBooking1To2AsUser2)
	t.Run("DeleteBookingHandler authorised", testDeleteBookingHandlerOfBooking1To2AsUser1)
	t.Run("GetBookingProvider after delete", testGetBookingProviderHandlerOfBooking1To2AfterDeletion)
}

// Helper
// ProviderSetting only makes providers on display on GET, booking can still be done for providers without settings
func testAddBookingHandlerToNoProviderSettingUser1ByUser0(t *testing.T) {
	ioReaderBooking, err := test_helper.GetIOReaderFromObject(test_helper.GetTestBooking(0, testUsers[1].ID))
	req, _ := http.NewRequest("POST", "/booking", ioReaderBooking)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to AddBooking failed with status code of %d", w.Code)
	}
	testBooking0to1, err = test_helper.GetBookingFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving new booking from response. %v", err)
	}
	if testBooking0to1.RecipientID != testUsers[0].ID {
		t.Errorf("Returned addedBooking did not update one of its RecipientID correctly")
	}
	if testBooking0to1.ProviderID != testUsers[1].ID {
		t.Errorf("Returned addedBooking did not update one of its ProviderID correctly")
	}
}

func testGetBookingProviderHandlerAsNotLoggedInNoProviderSetting(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/booking/%d", testBooking0to1.ID), nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusNotFound { 
		t.Errorf("HTTP Request to GetBookingProvider did not give status notfound but status code of %d", w.Code)
	}
}

func testGetAllBookingUsersHandlerAsNotLoggedIn(t *testing.T) {
	req, _ := http.NewRequest("GET", "/booking", nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllBookingUser failed with status code of %d", w.Code)
	}
	bookingUsers, err := test_helper.GetBookingUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all bookingUsers from response. %v", err)
	}
	if len(bookingUsers) != 0 {
		t.Errorf("A user who is not logged in saw some bookings involving him")
	}
}

func testGetAllBookingUserHandlerAsUser1(t *testing.T) {
	req, _ := http.NewRequest("GET", "/booking", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllBookingUser failed with status code of %d", w.Code)
	}
	bookingUsers, err := test_helper.GetBookingUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all bookingUsers from response. %v", err)
	}
	if len(bookingUsers) != 2 {
		t.Errorf("User1 does not see 2 bookings relevant to him")
	}
	if !test_helper.CheckBookingsInBookingUsers(bookingUsers, []Booking{ testBooking0to1, testBookingsTo2[1] }) {
		t.Errorf("Either booking0to1 or booking1to2 not found")
	}
}

func testGetAllBookingUsersHandlerSentAsUser1(t *testing.T) {
	req, _ := http.NewRequest("GET", "/booking?booking=SENT", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllBookingusers sent failed with status code of %d", w.Code)
	}
	bookingUsers, err := test_helper.GetBookingUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all bookingUsers from response. %v", err)
	}
	if len(bookingUsers) != 1 {
		t.Errorf("User1 did not see 1 booking sent by him")
	}
	if !test_helper.CheckBookingsInBookingUsers(bookingUsers, []Booking{ testBookingsTo2[1] }) {
		t.Errorf("booking1to2 not found")
	}
}

func testGetAllBookingUsersHandlerReceivedAsUser1(t *testing.T) {
	req, _ := http.NewRequest("GET", "/booking?booking=RECEIVED", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllBookingusers sent failed with status code of %d", w.Code)
	}
	bookingUsers, err := test_helper.GetBookingUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all bookingUsers from response. %v", err)
	}
	if len(bookingUsers) != 1 {
		t.Errorf("User1 did not see 1 booking received by him")
	}
	if !test_helper.CheckBookingsInBookingUsers(bookingUsers, []Booking{ testBooking0to1 }) {
		t.Errorf("booking0to1 not found")
	}
}

func testGetAllBookingUsersHandlerRequiredAsUser1(t *testing.T) {
	req, _ := http.NewRequest("GET", "/booking?booking=REQUIRED", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllBookingusers required failed with status code of %d", w.Code)
	}
	bookingUsers, err := test_helper.GetBookingUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all bookingUsers from response. %v", err)
	}
	if len(bookingUsers) != 1 {
		t.Errorf("User1 did not see 1 booking required by him")
	}
	if !test_helper.CheckBookingsInBookingUsers(bookingUsers, []Booking{ testBooking0to1 }) {
		t.Errorf("booking0to1 not found")
	}
}

func testGetAllBookingUsersHandlerRequiredAsUser2(t *testing.T) {
	req, _ := http.NewRequest("GET", "/booking?booking=REQUIRED", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[2],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllBookingUsers required failed with status code of %d", w.Code)
	}
	bookingUsers, err := test_helper.GetBookingUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all bookingUsers required from response. %v", err)
	}
	if len(bookingUsers) != 2 {
		t.Errorf("User2 does not see 2 bookingUsers required by him")
	}
	if !test_helper.CheckBookingsInBookingUsers(bookingUsers, testBookingsTo2) {
		t.Errorf("booking0to2 and booking1to2 not found")
	}
}

func testRespondBookingHandlerRejectNotLoggedIn(t *testing.T) {
	respond := BookingRespond{ Approve: false, Booking: Booking{ Details: "UpdatedDetails" } }
	ioReaderRespond, _ := test_helper.GetIOReaderFromObject(respond)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/booking/%d", testBooking0to1.ID), ioReaderRespond)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to respond while not logged in did not give status unauthorized but Status code: %d", w.Code)
	}
}

func testRespondBookingHandlerRejectAsUser1(t *testing.T) {
	respond := BookingRespond{ Approve: false, Booking: Booking{ Details: "UpdatedDetails" } }
	ioReaderRespond, _ := test_helper.GetIOReaderFromObject(respond)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/booking/%d", testBooking0to1.ID), ioReaderRespond)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to respond while authorized gave Status code: %d", w.Code)
	}
	respond, err := test_helper.GetBookingRespondFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving join request from response, %v", err)
	}
	testBooking0to1 = respond.Booking
}

func testGetAllBookingUserHandlerRequiredAsUser0(t *testing.T) {
	req, _ := http.NewRequest("GET", "/booking?booking=REQUIRED", nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK { 
		t.Errorf("HTTP Request to GetAllBookingUsers required failed with status code of %d", w.Code)
	}
	bookingUsers, err := test_helper.GetBookingUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving all bookingUsers required from response. %v", err)
	}
	if len(bookingUsers) != 1 {
		t.Errorf("User2 does not see 2 bookingUsers required by him")
	}
	if !bookingUsers[0].Booking.Equal(testBooking0to1) {
		t.Errorf("the bookingUser retrieved did not match updated booking")
	}
}

func testRespondBookingHandlerApproveAsUser0(t *testing.T) {
	bookingRespond := BookingRespond{ Approve: true }
	ioReaderRespond, _ := test_helper.GetIOReaderFromObject(bookingRespond)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/booking/%d", testBooking0to1.ID), ioReaderRespond)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[0],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to respond while authorized gave Status code: %d", w.Code)
	}
	eventWithUsers, err := test_helper.GetEventWithUsersFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving eventWithUsers from response, %v", err)
	}

	//Assert booking deleted
	rows, err := DB.Query(
		`SELECT COUNT(*) FROM wn_booking
		WHERE id = $1`,
		testBooking0to1.ID)
	if err != nil {
		t.Errorf("An error occured while getting count from DB. %v", err)
	}
	c, err := test_helper.ReadInt(rows)
	if err != nil {
		t.Errorf("An error occured while reading int. %v", err)
	}
	if c != 0 {
		t.Errorf("The booking still exist and has not been deleted")
	}

	//Assert event is created
	event := eventWithUsers.Event
	if event.ID == 0 {
		t.Errorf("Event was not properly created as it has an ID of 0")
	}
	if event.Access != "PRIVATE" {
		t.Errorf("The counsel event is not private but %s", event.Access)
	} 
	if l := len(eventWithUsers.Users); l != 2 {
		t.Errorf("The counsel event does not have 2 members but %d members", l)
	}
}

func testUpdateBookingHandlerOfBooking1To2NotLoggedIn(t *testing.T) {
	updatedBooking := Booking{ Nickname: "Simone Carter" }
	ioReaderBooking, _ := test_helper.GetIOReaderFromObject(updatedBooking)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/booking/%d", testBookingsTo2[1].ID), ioReaderBooking)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to update booking did not give unauthorised but gave Status code: %d", w.Code)
	}
}

func testUpdateBookingHandlerOfBooking1To2AsUser2Unauthorized(t *testing.T) {
	updatedBooking := Booking{ Nickname: "Simone Carter" }
	ioReaderBooking, _ := test_helper.GetIOReaderFromObject(updatedBooking)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/booking/%d", testBookingsTo2[1].ID), ioReaderBooking)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[2],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request to update booking did not give unauthorised but gave Status code: %d", w.Code)
	}
}

func testUpdateBookingHandlerOfBooking1To2AsUser1Authorized(t *testing.T) {
	updatedBooking := Booking{ Nickname: "Simone Carter" }
	ioReaderBooking, _ := test_helper.GetIOReaderFromObject(updatedBooking)
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/booking/%d", testBookingsTo2[1].ID), ioReaderBooking)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to update booking failed with Status code: %d", w.Code)
	}
	var err error
	testBookingsTo2[1], err = test_helper.GetBookingFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving Booking from response, %v", err)
	}
}

func testGetBookingProviderOfBooking1To2AfterUpdate(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/booking/%d", testBookingsTo2[1].ID), nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request to get booking failed with Status code: %d", w.Code)
	}
	bookingProvider, err := test_helper.GetBookingProviderFromRecorder(w)
	if err != nil {
		t.Errorf("An error occured while retrieving Booking from response, %v", err)
	}
	if bookingProvider.Booking.Nickname != "Simone Carter" {
		t.Errorf("booking1to2 is not updated properly")
	}
}

func testDeleteBookingHandlerOfBooking1To2AsUser2(t *testing.T) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/booking/%d", testBookingsTo2[1].ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[2],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP Request did not respond with unauthorized code but gave Status code: %d", w.Code)
	}
}

func testDeleteBookingHandlerOfBooking1To2AsUser1(t *testing.T) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/booking/%d", testBookingsTo2[1].ID), nil)
	req.AddCookie(&http.Cookie{
		Name: "session_key",
		Value: sessionKeys[1],
	})
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusOK {
		t.Errorf("HTTP Request did not respond with unauthorized code but gave Status code: %d", w.Code)
	}
}

func testGetBookingProviderHandlerOfBooking1To2AfterDeletion(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/booking/%d", testBookingsTo2[1].ID), nil)
	w := test_helper.SimulateRequest(Router, req)
	if w.Code != http.StatusNotFound { 
		t.Errorf("HTTP Request to GetBookingProvider did not respond with NotFound Code but with status code of %d", w.Code)
	}
}