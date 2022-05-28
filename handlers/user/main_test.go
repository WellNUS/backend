package user

// Test should be performed with some users in the database

import (
	"testing"
	"os"

	"fmt"
	"log"

	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

var (
	DOMAIN string = "localhost"
	FRONTEND_URL string = "localhost:3000"
	BACKEND_URL string = "localhost:8080"
	
	// Database fields
	HOST string = "localhost"
	PORT int = 5432
	USER string = "wellnus_user"
	PASSWORD string = "password"
	DB_NAME string = "wellnus"
)

func ConnectDB() *sql.DB {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
					USER,
					PASSWORD, 
					HOST,
					PORT,
					DB_NAME)
	// fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
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

func TestMain(m *testing.M) {
	db = ConnectDB()
	os.Exit(m.Run())
}