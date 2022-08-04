package testing

import (
	"wellnus/backend/config"
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/unit_test/test_helper"

	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTestingHomeHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		sID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		c.HTML(http.StatusOK, "home.html", gin.H{ "userID": sID, "backendURL": config.BACKEND_ADDRESS})
	}
}

func GetTestingAllUsersHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		users, _ := model.GetAllUsers(db)
		c.HTML(http.StatusOK, "users.html", gin.H{ "users": users })
	}
}

func GetTestingUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userIDCookie, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		userIDParam, _ := http_helper.GetIDParams(c)
		userWithGroups, _ := model.GetUserWithGroups(db, userIDParam)
		c.HTML(http.StatusOK, "user.html", gin.H{ 
			"userWithGroups": userWithGroups, 
			"backendURL": config.BACKEND_ADDRESS, 
			"wsURL": config.WS_ADDRESS,
			"userIDCookie": userIDCookie,
		})
	}
}

func GetTestingAllGroupsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		groups, _ := model.GetAllGroupsOfUser(db, userID)
		c.HTML(http.StatusOK, "groups.html", gin.H{ "groups": groups, "backendURL": config.BACKEND_ADDRESS })
	}
}

func GetTestingGroupHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		groupID, _ := http_helper.GetIDParams(c)
		groupWithUsers, _ := model.GetGroupWithUsers(db, groupID)
		c.HTML(http.StatusOK, "group.html", gin.H{"groupWithUsers": groupWithUsers, "backendURL": config.BACKEND_ADDRESS})
	}
}

func GetTestingAllJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		if s := c.Query("request"); s == "RECEIVED" {
			loadedJoinRequests, _ := model.GetAllLoadedJoinRequestsReceivedOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"loadedJoinRequests": loadedJoinRequests, "backendURL": config.BACKEND_ADDRESS})
		} else if s == "SENT" {
			loadedJoinRequests, _ := model.GetAllLoadedJoinRequestsSentOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"loadedJoinRequests": loadedJoinRequests, "backendURL": config.BACKEND_ADDRESS})
		} else {
			loadedJoinRequests, _ := model.GetAllLoadedJoinRequestsOfUser(db, userID)
			c.HTML(http.StatusOK, "joins.html", gin.H{"loadedJoinRequests": loadedJoinRequests, "backendURL": config.BACKEND_ADDRESS})
		}
	}
}

func GetTestingJoinRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		joinRequestID, _ := http_helper.GetIDParams(c)
		loadedJoinRequest, _ := model.GetLoadedJoinRequest(db, joinRequestID)
		c.HTML(http.StatusOK, "join.html", gin.H{"loadedJoinRequest": loadedJoinRequest, "backendURL": config.BACKEND_ADDRESS})
	}
}

func GetTestingChatHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		groupID, _ := http_helper.GetIDParams(c)
		groupWithUsers, _ := model.GetGroupWithUsers(db, groupID)
		c.HTML(http.StatusOK, "chat.html", gin.H{"groupWithUsers": groupWithUsers, "backendURL": config.BACKEND_ADDRESS, "wsURL": config.WS_ADDRESS})
	}
}

func GetTestingMatchHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		matchSetting, _ := model.GetMatchSettingOfUser(db, userID)
		count, _ := model.GetMatchRequestCount(db)
		c.HTML(http.StatusOK, "match.html", gin.H{"matchSetting": matchSetting, "mrCount": count, "backendURL": config.BACKEND_ADDRESS})
	}
}

func SetupUsersWithMatchRequests(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		count, err := strconv.Atoi(c.Query("count"))
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		users, err := test_helper.SetupUsers(db, count)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		_, err = test_helper.SetupMatchSettingForUsers(db, users)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		_, err =test_helper.SetupMatchRequestForUsers(db, users)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetTestingAllCounselRequestsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		topics, _ := c.GetQueryArray("topic")
		counselRequest, _ := model.GetCounselRequest(db, userID, userID)
		counselRequests, _ := model.GetAllCounselRequests(db, topics, userID)
		c.HTML(http.StatusOK, "counsel_requests.html", gin.H{"counselRequests": counselRequests, "counselRequest": counselRequest, "backendURL": config.BACKEND_ADDRESS})
	}
}

func GetTestingCounselRequestHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userIDParam, _ := http_helper.GetIDParams(c)
		userIDCookie, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		counselRequest, _ := model.GetCounselRequest(db, userIDParam, userIDCookie)
		c.HTML(http.StatusOK, "counsel_request.html", gin.H{"counselRequest": counselRequest, "backendURL": config.BACKEND_ADDRESS})
	}
}

func GetTestingAllEventsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		events, _ := model.GetAllEventsOfUser(db, userID)
		c.HTML(http.StatusOK, "events.html", gin.H{"events": events, "backendURL": config.BACKEND_ADDRESS})
	}
}

func GetTestingEventWithUsersHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		eventID, _ := http_helper.GetIDParams(c)
		eventWithUsers, _ := model.GetEventWithUsers(db, eventID)
		c.HTML(http.StatusOK, "event.html", gin.H{"eventWithUsers": eventWithUsers, "backendURL": config.BACKEND_ADDRESS})
	}
}

func GetTestingAllProvidersHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		topics, _ := c.GetQueryArray("topic")
		providerSetting, _ := model.GetProviderSetting(db, userID)
		providers, _ := model.GetAllProviders(db, topics)
		c.HTML(http.StatusOK, "providers.html", gin.H{"providers": providers, "providerSetting": providerSetting, "backendURL": config.BACKEND_ADDRESS})
	}
}

func GetTestingProviderWithEventsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userIDParam, _ := http_helper.GetIDParams(c)
		providerWithEvents, _ := model.GetProviderWithEvents(db, userIDParam)
		c.HTML(http.StatusOK, "provider.html", gin.H{"providerWithEvents": providerWithEvents, "backendURL": config.BACKEND_ADDRESS})
	}
} 

func GetTestingAllBookingUsersHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, _ := http_helper.GetUserIDFromSessionCookie(db, c)
		if s := c.Query("booking"); s == "RECEIVED" {
			bookingUsers, _ := model.GetAllBookingUsersReceivedOfUser(db, userID)
			c.HTML(http.StatusOK, "bookings.html", gin.H{"bookingUsers": bookingUsers, "backendURL": config.BACKEND_ADDRESS})
		} else if s == "SENT" {
			bookingUsers, _ := model.GetAllBookingUsersSentOfUser(db, userID)
			c.HTML(http.StatusOK, "bookings.html", gin.H{"bookingUsers": bookingUsers, "backendURL": config.BACKEND_ADDRESS})
		} else if s == "REQUIRED" {
			bookingUsers, _ := model.GetAllBookingUsersRequiredOfUser(db, userID)
			c.HTML(http.StatusOK, "bookings.html", gin.H{"bookingUsers": bookingUsers, "backendURL": config.BACKEND_ADDRESS})
		} else {
			bookingUsers, _ := model.GetAllBookingUsersOfUser(db, userID)
			c.HTML(http.StatusOK, "bookings.html", gin.H{"bookingUsers": bookingUsers, "backendURL": config.BACKEND_ADDRESS})
		}
	}
}

func GetTestingBookingProviderHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		bookingIDParam, _ := http_helper.GetIDParams(c)
		bookingProvider, _ := model.GetBookingProvider(db, bookingIDParam)
		c.HTML(http.StatusOK, "booking.html", gin.H{"bookingProvider": bookingProvider, "backendURL": config.BACKEND_ADDRESS})
	}
}