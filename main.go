package main

import (
	"wellnus/backend/config"

	"wellnus/backend/db"
	"wellnus/backend/router"
	"wellnus/backend/db/model"
)

func main() {
	config.LoadENV("./.env")

	// Runtime global instances
	DB := db.ConnectDB()
	WSHub := model.NewHub(DB)

	go WSHub.Run()
	Router := router.SetupRouter(DB, WSHub)

	Router.Run(config.SERVER_ADDRESS)
}