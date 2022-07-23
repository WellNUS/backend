package model

import (
	"time"
	"database/sql"
)

type Booking struct {
	ID			int64		`json:"id"`
	RecipientID int64 		`json:"recipient_id"`
	ProviderID 	int64		`json:"provider_id"`
	ApproveBy	int64		`json:"approve_by"`
	Nickname	string		`json:"nickname"`
	Details		string 		`json:"details"`
	StartTime	time.Time	`json:"start_time"`
	EndTime		time.Time	`json:"end_time"`
}

type BookingUser struct {
	Booking		Booking 	`json:"booking"`
	User		User		`json:"user"`
}

type BookingProvider struct {
	Booking 	Booking		`json:"booking"`
	Provider	Provider	`json:"provider"`		
}

type BookingRespond struct {
	Approve 	bool		`json:"approve"`
	Booking		Booking		`json:"booking"`
}

func (bMain Booking) MergeBooking(bAdd Booking) Booking {
	bMain.ID = bAdd.ID
	if bMain.RecipientID == 0 {
		bMain.RecipientID = bAdd.RecipientID
	}
	if bMain.ProviderID == 0 {
		bMain.ProviderID = bAdd.ProviderID
	}
	if bMain.ApproveBy == 0 {
		bMain.ApproveBy = bAdd.ApproveBy
	}
	if bMain.Nickname == "" {
		bMain.Nickname = bAdd.Nickname
	}
	if bMain.Details == "" {
		bMain.Details = bAdd.Details
	}
	if bMain.StartTime.IsZero() {
		bMain.StartTime = bAdd.StartTime
	}
	if bMain.EndTime.IsZero() {
		bMain.EndTime = bAdd.EndTime
	}
	return bMain
}

func (b Booking) LoadLastBookingID(db *sql.DB) (Booking, error) {
	row, err := db.Query("SELECT last_value FROM wn_booking_id_seq;")
	if err != nil { return Booking{}, err }
	defer row.Close()
	row.Next()
	if err := row.Scan(&b.ID); err != nil { return Booking{}, err }
	return b, nil
}

func (b Booking) LoadBookingWithProvider(db *sql.DB) (BookingProvider, error) {
	provider, err := GetProvider(db, b.ProviderID)
	if err != nil { return BookingProvider{}, err }
	return BookingProvider{ Booking: b, Provider: provider }, nil
}

func (b Booking) LoadBookingWithUser(db *sql.DB) (BookingUser, error) {
	user, err := GetUser(db, b.ProviderID)
	if err != nil { return BookingUser{}, err }
	return BookingUser{ Booking: b, User: user }, nil
}

func (b Booking) FlippedApproveBy() int64 {
	if b.ApproveBy == b.RecipientID {
		return b.ProviderID
	}
	return b.RecipientID
}

func (b1 Booking) Equal(b2 Booking) bool {
	return b1.ID == b2.ID &&
		b1.RecipientID == b2.RecipientID &&
		b1.ProviderID == b2.ProviderID &&
		b1.ApproveBy == b2.ApproveBy &&
		b1.Nickname == b2.Nickname &&
		b1.Details == b2.Details &&
		b1.StartTime.Equal(b2.StartTime) &&
		b1.EndTime.Equal(b2.EndTime)
}