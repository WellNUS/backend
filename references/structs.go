package references

type User struct {
	ID 				int64 	`json:"id"`
	FirstName 		string 	`json:"first_name"`
	LastName 		string	`json:"last_name"`
	Gender			string 	`json:"gender"`
	Email			string	`json:"email"`
	UserRole		string 	`json:"user_role"`
	Password		string 	`json:"password"`
	PasswordHash 	string	`json:"password_hash"`
}