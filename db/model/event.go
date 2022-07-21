package model

import (
	"time"
	"database/sql"
)

type UserIDBody struct {
	UserID	int64	`json:"user_id"`
}

type Event struct {
	ID 					int64 			`json:"id"`
	OwnerID 			int64 			`json:"owner_id"`
	EventName 			string 			`json:"event_name"`
	EventDescription 	string 			`json:"event_description"`
	StartTime 			time.Time		`json:"start_time"`
	EndTime				time.Time		`json:"end_time"`
	Access				string			`json:"access"`
	Category			string			`json:"category"`
}

type EventWithUsers struct {
	Event	Event	`json:"event"`
	Users	[]User	`json:"users"`
}

func (eventMain Event) MergeEvent(eventAdd Event) Event {
	eventMain.ID = eventAdd.ID
	if eventMain.OwnerID == 0 {
		eventMain.OwnerID = eventAdd.OwnerID
	}
	if eventMain.EventName == "" {
		eventMain.EventName = eventAdd.EventName
	}
	if eventMain.EventDescription == "" {
		eventMain.EventDescription = eventAdd.EventDescription
	}
	if eventMain.StartTime.IsZero() {
		eventMain.StartTime = eventAdd.StartTime
	}
	if eventMain.EndTime.IsZero() {
		eventMain.EndTime = eventAdd.EndTime
	}
	if eventMain.Access == "" {
		eventMain.Access = eventAdd.Access
	}
	if eventMain.Category == "" {
		eventMain.Category = eventAdd.Category
	}
	return eventMain
}

func (event Event) LoadLastEventID(db *sql.DB) (Event, error) {
	row, err := db.Query("SELECT last_value FROM wn_event_id_seq;")
	if err != nil { return Event{}, err }
	defer row.Close()
	row.Next()
	if err := row.Scan(&event.ID); err != nil { return Event{}, err }
	return event, nil
}

func (event1 Event) Equal(event2 Event) bool {
	return event1.ID == event2.ID &&
		event1.OwnerID == event2.OwnerID &&
		event1.EventName == event2.EventName &&
		event1.EventDescription == event2.EventDescription &&
		event1.StartTime == event2.StartTime &&
		event1.EndTime == event2.EndTime &&
		event1.Access == event2.Access &&
		event1.Category == event2.Category
}