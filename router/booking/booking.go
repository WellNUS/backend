package booking

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"

	"database/sql"
	"github.com/gin-gonic/gin"
)

const (
	BOOKING_RECEIVED 	= 0
	BOOKING_SENT 		= 1
	BOOKING_REQUIRED 	= 2
	BOOKING_ALL 		= 3
)

// Helper functions

func getBookingQuery(c *gin.Context) int {
	if s := c.Query("booking"); s == "RECEIVED" {
		return BOOKING_RECEIVED
	} else if s == "SENT" {
		return BOOKING_SENT
	} else if s == "REQUIRED" {
		return BOOKING_REQUIRED
	} else {
		return BOOKING_ALL
	}
}

// Main functions

func GetAllBookingUsersHandler(db *sql.DB) func(*gin.Context){
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, _ := http_helper.GetUserIDFromSessionHeader(db, c)
		request := getBookingQuery(c)
		if request == BOOKING_RECEIVED {
			bookingUsers, err := model.GetAllBookingUsersReceivedOfUser(db, userID)
			if err != nil {
				c.JSON(http_error.GetStatusCode(err), err.Error())
				return
			}
			c.JSON(http_error.GetStatusCode(err), bookingUsers)
		} else if request == BOOKING_SENT {
			bookingUsers, err := model.GetAllBookingUsersSentOfUser(db, userID)
			if err != nil {
				c.JSON(http_error.GetStatusCode(err), err.Error())
				return
			}
			c.JSON(http_error.GetStatusCode(err), bookingUsers)
		} else if request == BOOKING_REQUIRED {
			bookingUsers, err := model.GetAllBookingUsersRequiredOfUser(db, userID)
			if err != nil {
				c.JSON(http_error.GetStatusCode(err), err.Error())
				return
			}
			c.JSON(http_error.GetStatusCode(err), bookingUsers)
		} else {
			bookingUsers, err := model.GetAllBookingUsersOfUser(db, userID)
			if err != nil {
				c.JSON(http_error.GetStatusCode(err), err.Error())
				return
			}
			c.JSON(http_error.GetStatusCode(err), bookingUsers)
		}
	}
}

func GetBookingProviderHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		bookingIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		bookingProvider, err := model.GetBookingProvider(db, bookingIDParam)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), bookingProvider)
	}
}

func AddBookingHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetUserIDFromSessionHeader(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		booking, err := http_helper.GetBookingFromContext(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		booking, err = model.AddBooking(db, booking, booking.ProviderID, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), booking)
	}
}

func UpdateBookingHandler(db *sql.DB) func(*gin.Context) {	
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, _ := http_helper.GetUserIDFromSessionHeader(db, c)
		bookingIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		updatedBooking, err := http_helper.GetBookingFromContext(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		updatedBooking, err = model.UpdateBooking(db, updatedBooking, bookingIDParam, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), updatedBooking)
	}	
}

func RespondBookingHandler(db *sql.DB) func(*gin.Context) {	
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, _ := http_helper.GetUserIDFromSessionHeader(db, c)
		bookingIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		bookingRespond, err := http_helper.GetBookingRespondFromContext(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		// Either eventWithUsers or BookingRespond
		response, err := model.RespondBooking(db, bookingRespond, bookingIDParam, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), response)
	}	
}

func DeleteBookingHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetUserIDFromSessionHeader(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		bookingIDParam, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		booking, err := model.DeleteBookingAuthorized(db, bookingIDParam, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), booking)
	}
}