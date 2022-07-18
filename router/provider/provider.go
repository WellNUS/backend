package provider

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"

	"database/sql"
	"github.com/gin-gonic/gin"
)

func GetAllProvidersWithSettingHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		providersWithSetting, err := model.GetAllProvidersWithSetting(db)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), providersWithSetting)
	}
}

func GetProviderWithSettingHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetIDParams(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		providerWithSetting, err := model.GetProviderWithSetting(db, userID)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), providerWithSetting)
	}
}

func AddUpdateProviderSettingOfProviderHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		providerSetting, err := http_helper.GetProviderSettingFromContext(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		providerSetting, err = model.AddUpdateProviderSettingOfProvider(db, providerSetting, userID)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), providerSetting)
	}
}

func DeleteProviderSettingOfProviderHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		providerSetting, err := model.DeleteProviderSettingOfProvider(db, userID)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), providerSetting)
	}
}

