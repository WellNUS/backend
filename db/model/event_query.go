package model

import (
	"wellnus/backend/router/http_helper/http_error"	
	"log"
	"fmt"
	"database/sql"
	"errors"
)

// Helper function
func ReadEvents(rows *sql.Rows) ([]Event, error) {
	events := make([]Event, 0)
	for rows.Next() {
		var event Event
		if err := rows.Scan(
			&event.ID,
			&event.OwnerID,
			&event.EventName,
			&event.EventDescription,
			&event.StartTime,
			&event.EndTime,
			&event.Access); 
			err != nil {
				return nil, err
			}
		events = append(events, event)
	}
	return events, nil
}

func GetEvent(db *sql.DB, eventID int64) (Event, error) {
	rows, err := db.Query("SELECT * FROM wn_event WHERE id = $1;", eventID)
	if err != nil { return Event{}, err }
	defer rows.Close()
	events, err := ReadEvents(rows)
	if err != nil { return Event{}, err }
	if len(events) == 0 { return Event{}, http_error.NotFoundError }
	return events[0], nil
}

func AddUserToEvent(db *sql.DB, eventID int64, userID int64) error {
	_, err := db.Exec(
		`INSERT INTO wn_user_event (
			user_id, 
			event_id) 
		VALUES ($1, $2)`, 
		userID, 
		eventID)
	return err
}

func RemoveUserFromEvent(db *sql.DB, eventID int64, userID int64) error {
	_, err := db.Exec(
		`DELETE FROM wn_user_event WHERE
			user_id = $1 AND
			event_id = $2`,
		userID,
		eventID)
	return err
}

func DeleteEvent(db *sql.DB, eventID int64) error {
	_, err := db.Exec("DELETE FROM wn_event WHERE id = $1", eventID)
	return err
}

// Main Functions

func GetEventWithUsers(db *sql.DB, eventID int64) (EventWithUsers, error) {
	event, err := GetEvent(db, eventID)
	if err != nil { return EventWithUsers{}, err }
	users, err := GetAllUsersOfEvent(db, eventID)
	if err != nil { return EventWithUsers{}, err }
	return EventWithUsers{ Event: event, Users: users}, nil
}

func GetAllEventsOfUser(db *sql.DB, userID int64) ([]Event, error) {
	rows, err := db.Query(
		`SELECT
			wn_event.id,
			wn_event.owner_id,
			wn_event.event_name, 
			wn_event.event_description,
			wn_event.start_time, 
			wn_event.end_time,
			wn_event.access
		FROM wn_user_event JOIN wn_event 
		ON wn_user_event.event_id = wn_event.id 
		WHERE wn_user_event.user_id = $1`,
		userID)
	if err != nil { return nil, err }
	defer rows.Close()
	events, err := ReadEvents(rows)
	if err != nil { return nil, err}
	return events, nil
}

func AddEventWithUserIDs(db *sql.DB, event Event, userIDs []int64) (EventWithUsers, error) {
	if len(userIDs) == 0 { return EventWithUsers{}, errors.New("Insufficient users to form a event") }
	ownerID := userIDs[0]
	userIDs = userIDs[1:]
	event.OwnerID = ownerID //Taking first userID as ownerID
	_, err := db.Exec(
		`INSERT INTO wn_event (
			owner_id,
			event_name, 
			event_description,
			start_time, 
			end_time,
			access) 
		VALUES ($1, $2, $3, $4, $5, $6);`,
		event.OwnerID,
		event.EventName,
		event.EventDescription,
		event.StartTime,
		event.EndTime,
		event.Access)
	if err != nil { return EventWithUsers{}, err }
	event, err = event.LoadLastEventID(db)
	if err != nil { return EventWithUsers{}, err }

	// Adding Owner
	if err = AddUserToEvent(db, event.ID, ownerID); err != nil {
		log.Printf("Failed to add Owner: %v", err)
		if _, fatal := db.Exec("DELETE FROM wn_event WHERE id = $1", event.ID); fatal != nil {
			log.Fatal(fmt.Sprintf("Failed to remove added event after failing to add owner. Fatal: %v", fatal))
		}
		return EventWithUsers{}, err
	}

	// Adding Other Users
	for _, userID := range userIDs {
		if err := AddUserToEvent(db, event.ID, userID); err != nil {
			return EventWithUsers{}, err
		}
	}
	users, err := GetAllUsersOfEvent(db, event.ID)
	if err != nil { return EventWithUsers{}, err }

	return EventWithUsers{ Event: event, Users: users }, nil
}

func UpdateEvent(db *sql.DB, updatedEvent Event, eventID int64, userID int64) (Event, error) {
	targetEvent, err := GetEvent(db, eventID)
	if err != nil { return Event{}, err }
	if targetEvent.OwnerID != userID { return Event{}, http_error.UnauthorizedError }

	updatedEvent = updatedEvent.MergeEvent(targetEvent)

	_, err = db.Exec(
		`UPDATE wn_event SET 
			owner_id = $1,
			event_name = $2, 
			event_description = $3,
			start_time = $4, 
			end_time = $5,
			access = $6
		WHERE id = $7;`,
		updatedEvent.OwnerID,
		updatedEvent.EventName,
		updatedEvent.EventDescription,
		updatedEvent.StartTime,
		updatedEvent.EndTime,
		updatedEvent.Access,
		eventID)
	if err != nil { return Event{}, err }
	return updatedEvent, nil
}

func LeaveDeleteEvent(db *sql.DB, eventID int64, userID int64) (EventWithUsers, error) {
	targetEvent, err := GetEvent(db, eventID)
	if err != nil { return EventWithUsers{}, err }
	if targetEvent.OwnerID == userID {
		err = DeleteEvent(db, eventID)
		if err != nil { return EventWithUsers{}, err }
		return EventWithUsers{ Event: Event{ ID: eventID } }, nil
	} else {
		err = RemoveUserFromEvent(db, eventID, userID)
		if err != nil { return EventWithUsers{}, err }
		users, err := GetAllUsersOfEvent(db, eventID)
		if err != nil { return EventWithUsers{}, err }
		return EventWithUsers{ Event: targetEvent, Users: users }, nil
	}
}

func LeaveDeleteAllEvents(db *sql.DB, userID int64) ([]EventWithUsers, error) {
	events, err := GetAllEventsOfUser(db, userID)
	if err != nil { return nil, err}
	eventsWithUsers := make([]EventWithUsers, 0)
	for _, event := range events {
		eventWithUsers, err := LeaveDeleteEvent(db, event.ID, userID)
		if err != nil { return nil, err}
		eventsWithUsers = append(eventsWithUsers, eventWithUsers)
	}
	return eventsWithUsers, nil
}

func IsUserInEvent(db *sql.DB, userID int64, eventID int64) (bool, error) {
	row, err := db.Query(
		`SELECT COUNT(*) != 0 FROM wn_user_event 
		WHERE user_id = $1 and event_id = $2`,
		userID,
		eventID)
	if err != nil { return false, err }
	defer row.Close()
	var membership bool
	row.Next()
	if err := row.Scan(&membership); err != nil { return false, nil }
	return membership, nil
}

func AddUserToEventAuthorized(db *sql.DB, userID int64, eventID int64, adderID int64) (EventWithUsers, error) {
	targetEvent, err := GetEvent(db, eventID)
	if err != nil { return EventWithUsers{}, err }
	if targetEvent.Access == "PRIVATE" && targetEvent.OwnerID != adderID  {
		return EventWithUsers{}, http_error.UnauthorizedError
	}
	if err = AddUserToEvent(db, eventID, userID); err != nil {
		return EventWithUsers{}, err
	}
	users, err := GetAllUsersOfEvent(db, eventID)
	if err != nil { return EventWithUsers{}, err }
	return EventWithUsers{ Event: targetEvent, Users: users }, nil
}