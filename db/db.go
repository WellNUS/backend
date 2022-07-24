package db

import (
	"wellnus/backend/config"
	
	"fmt"
	"log"
	
	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("postgres", config.DB_ADDRESS)
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

func ConnectTestDB() *sql.DB {
	db, err := sql.Open("postgres", config.DB_ADDRESS_TEST)
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Database For Test Connected!")
	return db
}