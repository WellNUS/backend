package provider

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"

	"database/sql"
	"github.com/gin-gonic/gin"
)

func GetAllProvidersHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		topics, _ := c.GetQueryArray("topic")
		providers, err := model.GetAllProviders(db, topics)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		
		c.JSON(http_error.GetStatusCode(err), providers)
	}
}

func GetProviderWithEventsHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetIDParams(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		providerWithEvents, err := model.GetProviderWithEvents(db, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), providerWithEvents)
	}
}

func AddUpdateProviderSettingOfUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetUserIDFromSessionHeader(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		providerSetting, err := http_helper.GetProviderSettingFromContext(c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		providerSetting, err = model.AddUpdateProviderSettingOfUser(db, providerSetting, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), providerSetting)
	}
}

func DeleteProviderSettingOfUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetUserIDFromSessionHeader(db, c)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		providerSetting, err := model.DeleteProviderSettingOfUser(db, userID)
		if err != nil {
			c.JSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.JSON(http_error.GetStatusCode(err), providerSetting)
	}
}

