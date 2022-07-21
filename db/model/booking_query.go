package model

import (
	"wellnus/backend/router/http_helper/http_error"
	"database/sql"
	"fmt"
)

// Helper function

func ReadBookings(rows *sql.Rows) ([]Booking, error) {
	bookings := make([]Booking, 0)
	for rows.Next() {
		var booking Booking
		if err := rows.Scan(
			&booking.ID,
			&booking.RecipientID,
			&booking.ProviderID,
			&booking.ApproveBy,
			&booking.Nickname,
			&booking.Details,
			&booking.StartTime,
			&booking.EndTime); 
			err != nil {
				return nil, err
			}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func ReadBookingUsers(rows *sql.Rows) ([]BookingUser, error) {
	bookingUsers := make([]BookingUser, 0)
	for rows.Next() {
		var bookingUser BookingUser
		if err := rows.Scan(
			&bookingUser.Booking.ID,
			&bookingUser.Booking.RecipientID,
			&bookingUser.Booking.ProviderID,
			&bookingUser.Booking.ApproveBy,
			&bookingUser.Booking.Nickname,
			&bookingUser.Booking.Details,
			&bookingUser.Booking.StartTime,
			&bookingUser.Booking.EndTime,
			&bookingUser.User.ID, 
			&bookingUser.User.FirstName,
		 	&bookingUser.User.LastName, 
			&bookingUser.User.Gender, 
			&bookingUser.User.Faculty, 
			&bookingUser.User.Email, 
			&bookingUser.User.UserRole, 
			&bookingUser.User.PasswordHash); 
			err != nil {
				return nil, err
			}
		bookingUsers = append(bookingUsers, bookingUser)
	}
	return bookingUsers, nil
}

func GetBooking(db *sql.DB, bookingID int64) (Booking, error) {
	rows, err := db.Query("SELECT * FROM wn_counsel_booking WHERE id = $1", bookingID)
	if err != nil { return Booking{}, err }
	defer rows.Close()
	bookings, err := ReadBookings(rows)
	if err != nil { return Booking{}, err }
	if len(bookings) == 0 { return Booking{}, http_error.NotFoundError }
	return bookings[0], nil
}

func DeleteBooking(db *sql.DB, bookingID int64) (Booking, error) {
	_, err := db.Exec("DELETE FROM wn_counsel_booking WHERE id = $1", bookingID)
	if err != nil { return Booking{}, err }
	return Booking{ ID : bookingID }, nil
}

// Main function

func GetAllBookingUsersSentOfUser(db *sql.DB, userID int64) ([]BookingUser, error) {
	rows, err := db.Query(
		`SELECT 
			wn_counsel_booking.id, 
			wn_counsel_booking.recipient_id, 
			wn_counsel_booking.provider_id,
			wn_counsel_booking.approve_by,
			wn_counsel_booking.nickname,
			wn_counsel_booking.details,
			wn_counsel_booking.start_time,
			wn_counsel_booking.end_time,
			wn_user.id,
			wn_user.first_name,
			wn_user.last_name,
			wn_user.gender,
			wn_user.faculty,
			wn_user.email,
			wn_user.user_role,
			wn_user.password_hash
		FROM wn_counsel_booking JOIN wn_user
		ON wn_counsel_booking.provider_id = wn_user.id
		WHERE wn_counsel_booking.recipient_id = $1`,
		userID)
	if err != nil { return nil, err }
	defer rows.Close()
	bookingUsers, err := ReadBookingUsers(rows)
	if err != nil { return nil, err }
	return bookingUsers, nil
}

func GetAllBookingUsersReceivedOfUser(db *sql.DB, userID int64) ([]BookingUser, error) {
	rows, err := db.Query(
		`SELECT 
			wn_counsel_booking.id, 
			wn_counsel_booking.recipient_id, 
			wn_counsel_booking.provider_id,
			wn_counsel_booking.approve_by,
			wn_counsel_booking.nickname,
			wn_counsel_booking.details,
			wn_counsel_booking.start_time,
			wn_counsel_booking.end_time,
			wn_user.id,
			wn_user.first_name,
			wn_user.last_name,
			wn_user.gender,
			wn_user.faculty,
			wn_user.email,
			wn_user.user_role,
			wn_user.password_hash
		FROM wn_counsel_booking JOIN wn_user
		ON wn_counsel_booking.provider_id = wn_user.id
		WHERE wn_counsel_booking.provider_id = $1`,
		userID)
	if err != nil { return nil, err }
	defer rows.Close()
	bookingUsers, err := ReadBookingUsers(rows)
	if err != nil { return nil, err }
	return bookingUsers, nil
}

func GetAllBookingUsersRequiredOfUser(db *sql.DB, userID int64) ([]BookingUser, error) {
	rows, err := db.Query(
		`SELECT 
			wn_counsel_booking.id, 
			wn_counsel_booking.recipient_id, 
			wn_counsel_booking.provider_id,
			wn_counsel_booking.approve_by,
			wn_counsel_booking.nickname,
			wn_counsel_booking.details,
			wn_counsel_booking.start_time,
			wn_counsel_booking.end_time,
			wn_user.id,
			wn_user.first_name,
			wn_user.last_name,
			wn_user.gender,
			wn_user.faculty,
			wn_user.email,
			wn_user.user_role,
			wn_user.password_hash
		FROM wn_counsel_booking JOIN wn_user
		ON wn_counsel_booking.provider_id = wn_user.id
		WHERE wn_counsel_booking.approve_by = $2`,
		userID)
	if err != nil { return nil, err }
	defer rows.Close()
	bookingUsers, err := ReadBookingUsers(rows)
	if err != nil { return nil, err }
	return bookingUsers, nil
}

func GetAllBookingUsersOfUser(db *sql.DB, userID int64) ([]BookingUser, error) {
	rows, err := db.Query(
		`SELECT 
			wn_counsel_booking.id, 
			wn_counsel_booking.recipient_id, 
			wn_counsel_booking.provider_id,
			wn_counsel_booking.approve_by,
			wn_counsel_booking.nickname,
			wn_counsel_booking.details,
			wn_counsel_booking.start_time,
			wn_counsel_booking.end_time,
			wn_user.id,
			wn_user.first_name,
			wn_user.last_name,
			wn_user.gender,
			wn_user.faculty,
			wn_user.email,
			wn_user.user_role,
			wn_user.password_hash
		FROM wn_counsel_booking JOIN wn_user
		ON wn_counsel_booking.provider_id = wn_user.id
		WHERE wn_counsel_booking.recipient_id = $1 
		OR wn_counsel_booking.provider_id = $2`,
		userID,
		userID)
	if err != nil { return nil, err }
	defer rows.Close()
	bookingUsers, err := ReadBookingUsers(rows)
	if err != nil { return nil, err }
	return bookingUsers, nil
}

func GetBookingUser(db *sql.DB, bookingID int64) (BookingUser, error) {
	booking, err := GetBooking(db, bookingID)
	if err != nil { return BookingUser{}, err }
	bookingUser, err := booking.LoadBookingWithUser(db)
	if err != nil { return BookingUser{}, err }
	return bookingUser, nil
}

func GetBookingProvider(db *sql.DB, bookingID int64) (BookingProvider, error) {
	booking, err := GetBooking(db, bookingID)
	if err != nil { return BookingProvider{}, err }
	bookingProvider, err := booking.LoadBookingWithProvider(db)
	if err != nil { return BookingProvider{}, err }
	return bookingProvider, nil
}

func AddBooking(db *sql.DB, booking Booking, providerID int64, recipientID int64) (Booking, error) {
	if !AuthoriseProvider(db, providerID) { return Booking{}, http_error.UnauthorizedError }
	booking.RecipientID = recipientID
	booking.ProviderID = providerID
	_, err := db.Exec(
		`INSERT INTO wn_counsel_booking (
			recipient_id, 
			provider_id,
			approve_by,
			nickname,
			details,
			start_time,
			end_time,
		) values ($1, $2, $3, $4, $5, $6);`, 
		booking.RecipientID,
		booking.ProviderID,
		booking.ProviderID,
		booking.Nickname,
		booking.Details,
		booking.StartTime,
		booking.EndTime)
	if err != nil { return Booking{}, err }
	booking, err = booking.LoadLastBookingID(db)
	if err != nil { return Booking{}, err }
	return booking, nil
}

func UpdateBooking(db *sql.DB, updatedBooking Booking, bookingID int64, userID int64) (Booking, error) {
	targetBooking, err := GetBooking(db, bookingID)
	if err != nil { return Booking{}, err }
	if userID != targetBooking.RecipientID && userID != targetBooking.ApproveBy {
		return Booking{}, http_error.UnauthorizedError
	}
	updatedBooking = updatedBooking.MergeBooking(targetBooking)
	_, err = db.Exec(
		`UPDATE wn_counsel_booking SET 
			recipient_id = $1, 
			provider_id = $2,
			approve_by = $3,
			nickname = $4,
			details = $5,
			start_time = $6,
			end_time = $7,
		WHERE id = $8;`,
		updatedBooking.RecipientID,
		updatedBooking.ProviderID,
		updatedBooking.ApproveBy,
		updatedBooking.Nickname,
		updatedBooking.Details,
		updatedBooking.StartTime,
		updatedBooking.EndTime)
	if err != nil { return Booking{}, err }
	return updatedBooking, nil
}

func RespondBooking(db *sql.DB, bookingRespond BookingRespond, bookingID int64, userID int64) (BookingRespond, error) {
	bookingUser, err := GetBookingUser(db, bookingID)
	if err != nil { return BookingRespond{}, nil }
	if bookingUser.Booking.ApproveBy != userID {
		return BookingRespond{}, http_error.UnauthorizedError 
	}

	booking := bookingUser.Booking
	user := bookingUser.User
	if bookingRespond.Approve {
		nickname := booking.Nickname
		providerName := user.FirstName
		event := Event{
			EventName: fmt.Sprintf("%s and %s Counsel Session", nickname, providerName),
			EventDescription: fmt.Sprint("Counsel Session for %s by %s", nickname, providerName),
			StartTime: booking.StartTime,
			EndTime: booking.EndTime,
			Access: "PRIVATE",
			Category: "COUNSEL",
		}
		_, err = AddEventWithUserIDs(db, event, []int64{booking.ProviderID, booking.RecipientID})
		if err != nil { return BookingRespond{}, err }
		_, err = DeleteBooking(db, bookingID)
		if err != nil { return BookingRespond{}, err }
	} else {
		updatedBooking := bookingRespond.Booking.MergeBooking(booking)
		updatedBooking.ApproveBy = booking.FlippedApproveBy()
		_, err = UpdateBooking(db, updatedBooking, bookingID, booking.RecipientID)
		if err != nil { return BookingRespond{}, err }
	}
	return bookingRespond, nil
}

func DeleteBookingAuthorized(db *sql.DB, bookingID int64, userID int64) (Booking, error) {
	targetBooking, err := GetBooking(db, bookingID)
	if err != nil { return Booking{}, err }
	if targetBooking.RecipientID != userID { return Booking{}, err }
	booking, err := DeleteBooking(db, bookingID)
	if err != nil { return Booking{}, err }
	return booking, nil
}

