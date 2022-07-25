package db

import (
	"wellnus/backend/config"
	
	"fmt"
	"log"
	
	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	address := config.DB_ADDRESS
	if config.RUN_WITH_DOCKER_COMPOSE {
		fmt.Println("Database starting with docker compose")
		address = config.DOCKER_COMPOSE_DB_ADDRESS
	}
	
	db, err := sql.Open("postgres", address)
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Database Connected!")
	return db
}