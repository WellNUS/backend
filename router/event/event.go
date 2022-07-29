package event

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"

	"github.com/gin-gonic/gin"
	"database/sql"
)

// Main functions
func GetAllEventsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, _ := http_helper.GetUserIDFromSessionHeader(db, c)
		events, err := model.GetAllEventsOfUser(db, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), events)
	}
}

func GetEventHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		eventIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		eventWithUsers, err := model.GetEventWithUsers(db, eventIDParam)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), eventWithUsers)
	}
}

func AddEventHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		newEvent, err := http_helper.GetEventFromContext(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		userID, err := http_helper.GetUserIDFromSessionHeader(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		eventWithUsers, err := model.AddEventWithUserIDs(db, newEvent, []int64{userID}) // Can throw a fatal error
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), eventWithUsers)
	}
}

func UpdateEventHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		eventIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, err := http_helper.GetUserIDFromSessionHeader(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		updatedEvent, err := http_helper.GetEventFromContext(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		updatedEvent, err = model.UpdateEvent(db, updatedEvent, eventIDParam, userIDCookie)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), updatedEvent)
	}
}

func LeaveDeleteEventHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		eventIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userIDCookie, err := http_helper.GetUserIDFromSessionHeader(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		eventWithUsers, err := model.LeaveDeleteEvent(db, eventIDParam, userIDCookie)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), eventWithUsers)
	}
}

func LeaveDeleteAllEventsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userIDCookie, err := http_helper.GetUserIDFromSessionHeader(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		eventsWithUsers, err := model.LeaveDeleteAllEvents(db, userIDCookie)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), eventsWithUsers)
	}
}

func AddUserToEventHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userIDCookie, err := http_helper.GetUserIDFromSessionHeader(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		userIDAdded, err := http_helper.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		eventIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		eventWithUsers, err := model.AddUserToEventAuthorized(db, userIDAdded, eventIDParam, userIDCookie)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), eventWithUsers)
	}
}

func CreateGroupDeleteEventHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userIDCookie, err := http_helper.GetUserIDFromSessionHeader(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		eventIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		groupWithUsers, err := model.CreateGroupDeleteEvent(db, eventIDParam, userIDCookie)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), groupWithUsers)
	}
}