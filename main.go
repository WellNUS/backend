package main

import (
	"wellnus/backend/config"

	"wellnus/backend/db"
	"wellnus/backend/router"
	"wellnus/backend/router/ws"
)

func main() {
	config.LoadENV("./.env")

	// Runtime global instances
	DB := db.ConnectDB()
	WSHub := ws.NewHub(DB)

	go WSHub.Run()
	Router := router.SetupRouter(DB, WSHub)

	Router.Run(config.BACKEND_DOMAIN)
}