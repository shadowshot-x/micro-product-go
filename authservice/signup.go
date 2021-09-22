package authservice

import (
	"net/http"

	"github.com/shadowshot-x/micro-product-go/authservice/data"
)

// adds the user to the database of users
func SignupHandler(rw http.ResponseWriter, r *http.Request) {
	check := data.AddUserObject(r.Header["Email"][0], r.Header["Username"][0], r.Header["Passwordhash"][0],
		r.Header["Fullname"][0], 0)

	if !check {
		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte("Email or Username already exists"))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("User Created"))
}
