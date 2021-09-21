package authservice

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/shadowshot-x/micro-product-go/authservice/data"
)

// we need this function to be private
func getSignedToken() (string, error) {
	// we make a JWT Token here with signing method of ES256 and claims.
	// claims are attributes.
	// aud - audience
	// iss - issuer
	// exp - expiration of the Token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"aud": "frontend.knowsearch.ml",
		"iss": "knowsearch.ml",
		"exp": time.Now().Add(time.Minute * 1).Unix(),
	})
	// here we provide the shared secret. It should be very complex.\
	// Aslo, it should be passed as a System Environment variable
	tokenString, err := token.SignedString([]byte("S0m3_R7nd0m_ST41ng"))
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

// searches the user in the database.
func validateUser(email string, passwordHash string) (bool, error) {
	usr, exists := data.GetUserObject(email)
	if !exists {
		return false, errors.New("user does not exist")
	}
	passwordCheck := usr.ValidatePasswordHash(passwordHash)

	if !passwordCheck {
		return false, nil
	}
	return true, nil
}

// This will be supplied to the MUX router. It will be called when signin request is sent
// if user not found or not validates, returns the Unauthorized error
// if found, returns the JWT back. [How to return this?]
func SigninHandler(rw http.ResponseWriter, r *http.Request) {
	// lets see if the user exists
	valid, err := validateUser(r.Header["Email"][0], r.Header["Passwordhash"][0])
	if err != nil {
		// this means either the user does not exist
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("User Does not Exist"))
		return
	}

	if !valid {
		// this means either the password is wrong
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Incorrect Password"))
		return
	}
	tokenString, err := getSignedToken()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Internal Server Error"))
		return
	}

	fmt.Println(tokenString)
}
