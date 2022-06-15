package db

import (
	"wellnus/backend/config"
	
	"fmt"
	"log"
	
	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	// fmt.Println(connStr)
	db, err := sql.Open("postgres", config.CONNECTION_STRING)
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