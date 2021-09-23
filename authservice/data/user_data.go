package data

// User struct
type user struct {
	email        string
	username     string
	passwordhash string
	fullname     string
	createDate   string
	role         int
}

// currently this is acting as our database
var userList = []user{
	{
		email:        "abc@gmail.com",
		username:     "abc12",
		passwordhash: "hashedme1",
		fullname:     "abc def",
		createDate:   "1631600786",
		role:         1,
	},
	{
		email:        "chekme@example.com",
		username:     "checkme34",
		passwordhash: "hashedme2",
		fullname:     "check me",
		createDate:   "1631600837",
		role:         0,
	},
}

// based on the email id provided, finds the user object
// can be seen as the main constructor to start validation
func GetUserObject(email string) (user, bool) {
	//needs to be replaces using Database
	for _, user := range userList {
		if user.email == email {
			return user, true
		}
	}
	return user{}, false
}

// checks if the password hash is valid
func (u *user) ValidatePasswordHash(pswdhash string) bool {
	return u.passwordhash == pswdhash
}

// this simply adds the user to the list
func AddUserObject(email string, username string, passwordhash string, fullname string, role int) bool {
	// declare the new user object
	newUser := user{
		email:        email,
		passwordhash: passwordhash,
		username:     username,
		fullname:     fullname,
		role:         role,
	}
	// check if a user already exists
	for _, ele := range userList {
		if ele.email == email || ele.username == username {
			return false
		}
	}
	userList = append(userList, newUser)
	return true
}
