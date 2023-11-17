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
	fmt.Println("Connecting to database: ", address)
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