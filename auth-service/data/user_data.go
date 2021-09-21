package authdata

import "errors"

// User struct
type User struct {
	Email        string
	Username     string
	Passwordhash string
	Fullname     string
	CreateDate   string
}

// currently this is acting as our database
var UserList = []User{
	{
		Email:        "abc@gmail.com",
		Username:     "abc12",
		Passwordhash: "hashedme1",
		Fullname:     "abc def",
		CreateDate:   "1631600786",
	},
	{
		Email:        "chekme@example.com",
		Username:     "checkme34",
		Passwordhash: "hashedme2",
		Fullname:     "check me",
		CreateDate:   "1631600837",
	},
}

// based on the email id provided, finds the user object
// can be seen as the main constructor to start validation
func GetUserObject(email string) (User, error) {
	for _, user := range UserList {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, errors.New("User not found in the database")
}

// checks if the password hash is valid
func (u *User) ValidatePasswordHash(pswdhash string) bool {
	return u.Passwordhash == pswdhash
}
