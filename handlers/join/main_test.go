package join

import (
	"wellnus/backend/config"
	"wellnus/backend/handlers/misc"

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
	NotFoundErrorMessage 		string = misc.NotFoundError.Error()
	UnauthorizedErrorMessage	string = misc.UnauthorizedError.Error()
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

func equal_joinRequest(joinRequest1 JoinRequest, joinRequest2 JoinRequest) bool {
	return joinRequest1.ID == joinRequest2.ID &&
	joinRequest1.GroupID == joinRequest2.GroupID &&
	joinRequest1.UserID == joinRequest2.UserID
}

func equal_user(user1 User, user2 User) bool {
	return user1.ID == user2.ID &&
	user1.FirstName == user2.FirstName &&
	user1.LastName == user2.LastName &&
	user1.Gender == user2.Gender &&
	user1.Faculty == user2.Faculty &&
	user1.Email == user2.Email &&
	user1.UserRole == user2.UserRole &&
	user1.PasswordHash == user2.PasswordHash
}

func equal_group(group1 Group, group2 Group) bool {
	return group1.ID == group2.ID &&
	group1.GroupName == group2.GroupName &&
	group1.GroupDescription == group2.GroupDescription &&
	group1.Category == group2.Category &&
	group1.OwnerID == group2.OwnerID
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

func addUserToGroup(db *sql.DB, groupID int64, userID int64) error {
	_, err := db.Exec(
		`INSERT INTO wn_user_group (
			user_id, 
			group_id) 
		VALUES ($1, $2)`, 
		userID, 
		groupID)
	return err
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
	newUser, err = loadLastUserID(db, newUser)
	if err != nil { return User{}, err }
	return newUser, nil
}

func AddGroup(newGroup Group) (Group, error) {
	_, err := db.Query(
		`INSERT INTO wn_group (
			group_name, 
			group_description, 
			category, 
			owner_id) 
		VALUES ($1, $2, $3, $4);`,
		newGroup.GroupName,
		newGroup.GroupDescription,
		newGroup.Category,
		newGroup.OwnerID)
	if err != nil { return Group{}, err }
	newGroup, err = loadLastGroupID(db, newGroup)
	if err != nil { return Group{}, err }
	
	// newGroup successfully added into DB. Now adding owner into new group
	err = addUserToGroup(db, newGroup.ID, newGroup.OwnerID)
	if err != nil {
		log.Printf("Failed to add Owner: %v", err)
		if _, fatal := db.Exec("DELETE FROM wn_group WHERE id = $1`", newGroup.ID); fatal != nil {
			log.Fatal(fmt.Sprintf("Failed to remove added group after failing to add owner. Fatal: %v", fatal))
		}
		return Group{}, err
	}
	return newGroup, nil
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
	router.GET("/join", GetAllJoinRequestsHandler(db))
	router.POST("/join", AddJoinRequestHandler(db))
	router.GET("/join/:id", GetLoadedJoinRequestHandler(db))
	router.PATCH("/join/:id", RespondJoinRequestHandler(db))
	router.DELETE("/join/:id", DeleteJoinRequestHandler(db))

	fmt.Printf("Starting backend server at '%s' \n", config.BACKEND_URL)
	return router
}

func TestMain(m *testing.M) {
	db = setupDB()
	router = setupRouter()
	
	var err error
	validAddedUser1, err = makeNewUser(validAddedUser1)
	validAddedUser2, err = makeNewUser(validAddedUser2)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test user. %v", err)) }
	validAddedGroup.OwnerID = validAddedUser1.ID	//Setting user1 as owner
	validAddedGroup, err = AddGroup(validAddedGroup)
	if err != nil { log.Fatal(fmt.Sprintf("Something went wrong when creating Test group. %v", err)) }

	r := m.Run()

	db.Exec("DELETE FROM wn_group WHERE id = $1", validAddedGroup.ID)
	db.Exec("DELETE FROM wn_user WHERE id = $1 OR id = $2", validAddedUser1.ID, validAddedUser2.ID)
	os.Exit(r)
}