package session

import (
	"wellnus/backend/references"

	"testing"
	"os"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/alexedwards/argon2id"
	"database/sql"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	router *gin.Engine
	templateUser User = User{
		FirstName: "NewFirstName",
		LastName: "NewLastName",
		Gender: "M",
		Faculty: "COMPUTING",
		Email: "NewEmail@u.nus.edu",
		UserRole: "VOLUNTEER",
		Password: "NewPassword",
		PasswordHash: "",
	}
	
) 

func hashPassword(user User) (User, error) {
	var err error
	user.PasswordHash, err = argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	user.Password = ""
	if err != nil { return User{}, err }
	return user, nil
}

func loadLastID(db *sql.DB, user User) (User, error) {
	row, err := db.Query("SELECT last_value FROM wn_user_id_seq;")
	if err != nil { return User{}, err }
	defer row.Close()

	row.Next()
	if err := row.Scan(&user.ID); err != nil { return User{}, err }
	return user, nil
}

func makeNewUser(newUser User) (User, error) {
	newUser, err := hashPassword(newUser);
	if err != nil { return User{}, err }
	_, err = db.Query(fmt.Sprintf(
		"INSERT INTO wn_user (first_name, last_name, gender, faculty, email, user_role, password_hash) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');",
		newUser.FirstName,
		newUser.LastName,
		newUser.Gender,
		newUser.Faculty,
		newUser.Email,
		newUser.UserRole,
		newUser.PasswordHash))
	if err != nil { return User{}, err }
	// New user successfully made
	newUser, err = loadLastID(db, newUser)
	if err != nil { return User{}, err }
	return newUser, nil
}

func connectDB() *sql.DB {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
					references.USER,
					references.PASSWORD, 
					references.HOST,
					references.PORT,
					references.DB_NAME)
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

func setupRouter() *gin.Engine {
	router := gin.Default()
	
	router.POST("/session", LoginHandler(db))
	router.DELETE("/session", LogoutHandler(db))

	fmt.Printf("Starting backend server at '%s' \n", references.BACKEND_URL)
	return router
}

func TestMain(m *testing.M) {
	db = connectDB()
	router = setupRouter()
	user, err := makeNewUser(templateUser)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err)) }

	r := m.Run()

	_ , err = db.Query(fmt.Sprintf("DELETE FROM wn_user WHERE id = %d", user.ID))
	if err != nil { log.Fatal("Test user was not removed from database") }
	os.Exit(r)
}