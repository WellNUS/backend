package join

import (
	"wellnus/backend/references"
	"wellnus/backend/handlers/httpError"

	"testing"
	"os"
	"fmt"
	"log"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"database/sql"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB 
	router *gin.Engine
	NotFoundErrorMessage 		string = httpError.NotFoundError.Error()
	UnauthorizedErrorMessage	string = httpError.UnauthorizedError.Error()
)

var validAddedUser1 User = User{
	FirstName: "NewFirstName",
	LastName: "NewLastName",
	Gender: "M",
	Faculty: "COMPUTING",
	Email: "NewEmail@u.nus.edu",
	UserRole: "VOLUNTEER",
	Password: "NewPassword",
	PasswordHash: "",
}

var validAddedUser2 User = User{
	FirstName: "NewFirstName1",
	LastName: "NewLastName1",
	Gender: "M",
	Faculty: "COMPUTING",
	Email: "NewEmail1@u.nus.edu",
	UserRole: "VOLUNTEER",
	Password: "NewPassword",
	PasswordHash: "",
}

var validAddedGroup Group = Group{
	GroupName: "NewGroupName",
	GroupDescription: "NewGroupDescription",
	Category: "SUPPORT",
}


func equal(user1 User, user2 User) bool {
	return user1.ID == user2.ID &&
	user1.FirstName == user2.FirstName &&
	user1.LastName == user2.LastName &&
	user1.Gender == user2.Gender &&
	user1.Faculty == user2.Faculty &&
	user1.Email == user2.Email &&
	user1.UserRole == user2.UserRole &&
	user1.PasswordHash == user2.PasswordHash
}

func hashPassword(user User) (User, error) {
	var err error
	user.PasswordHash, err = argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	user.Password = ""
	if err != nil { return User{}, err }
	return user, nil
}

func loadLastUserID(db *sql.DB, user User) (User, error) {
	row, err := db.Query("SELECT last_value FROM wn_user_id_seq;")
	if err != nil { return User{}, err }
	defer row.Close()

	row.Next()
	if err := row.Scan(&user.ID); err != nil { return User{}, err }
	return user, nil
}

func loadLastGroupID(db *sql.DB, group Group) (Group, error) {
	row, err := db.Query("SELECT last_value FROM wn_group_id_seq;")
	if err != nil { return Group{}, err }
	defer row.Close()

	row.Next()
	if err := row.Scan(&group.ID); err != nil { return Group{}, err }
	return group, nil
}

func makeNewUser(newUser User) (User, error) {
	newUser, err := hashPassword(newUser);
	if err != nil { return User{}, err }
	_, err = db.Query(fmt.Sprintf(
		`INSERT INTO wn_user (
			first_name, 
			last_name, 
			gender, 
			faculty, 
			email, 
			user_role, 
			password_hash
		) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');`,
		newUser.FirstName,
		newUser.LastName,
		newUser.Gender,
		newUser.Faculty,
		newUser.Email,
		newUser.UserRole,
		newUser.PasswordHash))
	if err != nil { return User{}, err }
	// New user successfully made
	newUser, err = loadLastUserID(db, newUser)
	if err != nil { return User{}, err }
	return newUser, nil
}

func AddGroup(newGroup Group) (Group, error) {
	query := fmt.Sprintf(
		`INSERT INTO wn_group (
			group_name, 
			group_description, 
			category, 
			owner_id) 
		VALUES ('%s', '%s', '%s', %d);`,
		newGroup.GroupName,
		newGroup.GroupDescription,
		newGroup.Category,
		newGroup.OwnerID)
	_, err := db.Query(query)
	if err != nil { return Group{}, err }
	newGroup, err = loadLastGroupID(db, newGroup)
	if err != nil { return Group{}, err }
	
	// newGroup successfully added into DB. Now adding owner into new group
	err = addUserToGroup(db, newGroup.ID, newGroup.OwnerID)
	if err != nil {
		log.Printf("Failed to add Owner: %v", err)
		if _, fatal := db.Query(fmt.Sprintf("DELETE FROM wn_group WHERE id = %d", newGroup.ID)); fatal != nil {
			log.Fatal(fmt.Sprintf("Failed to remove added group after failing to add owner. Fatal: %v", fatal))
		}
		return Group{}, err
	}
	return newGroup, nil
}

func connectDB() *sql.DB {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
					references.USER,
					references.PASSWORD, 
					references.HOST,
					references.PORT,
					references.DB_NAME)
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

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/join", GetAllJoinRequestsHandler(db))
	router.POST("/join", AddJoinRequestHandler(db))
	router.GET("/join/:id", GetJoinRequestHandler(db))
	router.PATCH("/join/:id", RespondJoinRequestHandler(db))
	router.DELETE("/join/:id", DeleteJoinRequestHandler(db))

	fmt.Printf("Starting backend server at '%s' \n", references.BACKEND_URL)
	return router
}

func TestMain(m *testing.M) {
	db = connectDB()
	router = setupRouter()
	if _, err := db.Query("DELETE FROM wn_group;"); err != nil {
		log.Fatal(fmt.Sprintf("Unable to clear wn_group in preparation for test. %v", err))
	}
	if _, err := db.Query("DELETE FROM wn_user;"); err != nil {
		log.Fatal(fmt.Sprintf("Unable to clear wn_user in preparation for test. %v", err))
	}
	
	var err error
	validAddedUser1, err = makeNewUser(validAddedUser1)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err)) }
	validAddedUser2, err = makeNewUser(validAddedUser2)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err)) }
	validAddedGroup.OwnerID = validAddedUser1.ID	//Setting user1 as owner
	validAddedGroup, err = AddGroup(validAddedGroup)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test group. %v", err)) }

	r := m.Run()

	_ , err = db.Query(fmt.Sprintf("DELETE FROM wn_group WHERE id = %d", validAddedGroup.ID))
	if err != nil { log.Fatal("Test group was not removed from database") }
	_ , err = db.Query(fmt.Sprintf("DELETE FROM wn_user WHERE id = %d", validAddedUser1.ID))
	if err != nil { log.Fatal("Test user1 was not removed from database") }
	_ , err = db.Query(fmt.Sprintf("DELETE FROM wn_user WHERE id = %d", validAddedUser2.ID))
	if err != nil { log.Fatal("Test user2 was not removed from database") }
	
	os.Exit(r)
}