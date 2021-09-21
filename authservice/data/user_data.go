package data

// User struct
type User struct {
	email        string
	username     string
	passwordhash string
	fullname     string
	createDate   string
}

// currently this is acting as our database
var UserList = []User{
	{
		email:        "abc@gmail.com",
		username:     "abc12",
		passwordhash: "hashedme1",
		fullname:     "abc def",
		createDate:   "1631600786",
	},
	{
		email:        "chekme@example.com",
		username:     "checkme34",
		passwordhash: "hashedme2",
		fullname:     "check me",
		createDate:   "1631600837",
	},
}

// based on the email id provided, finds the user object
// can be seen as the main constructor to start validation
func GetUserObject(email string) (User, bool) {
	//needs to be replaces using Database
	for _, user := range UserList {
		if user.email == email {
			return user, true
		}
	}
	return User{}, false
}

// checks if the password hash is valid
func (u *User) ValidatePasswordHash(pswdhash string) bool {
	return u.passwordhash == pswdhash
}
