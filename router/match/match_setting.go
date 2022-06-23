package match

import (
	"wellnus/backend/db/model"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"

	"database/sql"
	"github.com/gin-gonic/gin"
)

func GetMatchSettingOfUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}

		matchSetting, err := model.GetMatchSettingOfUser(db, userID)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), matchSetting)
	}
}

func AddUpdateMatchSettingOfUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		matchSetting, err := http_helper.GetMatchSettingFromContext(c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		matchSetting, err = model.AddUpdateMatchSettingOfUser(db, matchSetting, userID)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), matchSetting)
	}
}

func DeleteMatchSettingOfUserHandler(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		matchSetting, err := model.DeleteMatchSettingOfUser(db, userID)
		if err != nil {
			c.IndentedJSON(http_error.GetStatusCode(err), err.Error())
			return
		}
		c.IndentedJSON(http_error.GetStatusCode(err), matchSetting)
	}
}