package session

import (
	"wellnus/backend/config"

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
) 

var validUser User = User{
		FirstName: "NewFirstName",
		LastName: "NewLastName",
		Gender: "M",
		Faculty: "COMPUTING",
		Email: "NewEmail@u.nus.edu",
		UserRole: "VOLUNTEER",
		Password: "NewPassword",
		PasswordHash: "",
	}

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
	_, err = db.Exec(
		`INSERT INTO wn_user (
			first_name, 
			last_name, 
			gender, 
			faculty, 
			email, 
			user_role, 
			password_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		newUser.FirstName,
		newUser.LastName,
		newUser.Gender,
		newUser.Faculty,
		newUser.Email,
		newUser.UserRole,
		newUser.PasswordHash)
	if err != nil { return User{}, err }
	// New user successfully made
	newUser, err = loadLastID(db, newUser)
	if err != nil { return User{}, err }
	return newUser, nil
}

func setupDB() *sql.DB {
	db, err := sql.Open("postgres", config.CONNECTION_STRING)
	if err != nil {
		log.Fatal(err)
	}
	db.Query("DELETE FROM wn_group;")
	db.Query("DELETE FROM wn_user;")
	return db
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	
	router.POST("/session", LoginHandler(db))
	router.DELETE("/session", LogoutHandler(db))

	fmt.Printf("Starting backend server at '%s' \n", config.BACKEND_URL)
	return router
}

func TestMain(m *testing.M) {
	db = setupDB()
	router = setupRouter()
	user, err := makeNewUser(validUser)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err)) }

	r := m.Run()

	db.Exec("DELETE FROM wn_user WHERE id = $1", user.ID)
	os.Exit(r)
}