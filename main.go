package main

import (
	"wellnus/backend/config"

	"wellnus/backend/db"
	"wellnus/backend/router"
	"wellnus/backend/router/ws"
	
	"fmt"
)

func main() {
	// Runtime global instances
	DB := db.ConnectDB()
	WSHub := ws.NewHub(DB)

	go WSHub.Run()
	Router := router.SetupRouter(DB, WSHub)

	fmt.Printf("Starting backend server at '%s' \n", config.BACKEND_URL)
	Router.Run(config.BACKEND_DOMAIN)
}

