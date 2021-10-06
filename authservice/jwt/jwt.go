package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	TOKEN_IS_CORRUPT = "Token is corrupt"
	INVALID_TOKEN    = "Invalid Token"
	TOKEN_EXPIRED    = "Token Expired"
)

// claims are attributes.
// Aud - audience
// Iss - issuer
// Exp - expiration of the Token
type ClaimsMap struct {
	Aud string
	Iss string
	Exp string
}

// GetSecret fetches the value for the JWT_SECRET from the environment variable
func GetSecret() string {
	return os.Getenv("JWT_SECRET")
}

// Function for generating the tokens.
func GenerateToken(header string, payload ClaimsMap, secret string) (string, error) {
	// create a new hash of type sha256. We pass the secret key to it
	// sha256 is a symmetric cryptographic algorithm
	h := hmac.New(sha256.New, []byte(secret))

	// We base encode the header which is a normal string
	header64 := base64.StdEncoding.EncodeToString([]byte(header))
	// We then Marshal the payload which is a map. This converts it to a string of JSON.
	// Now we base encode this string
	payloadstr, err := json.Marshal(payload)
	if err != nil {
		return string(payloadstr), fmt.Errorf("Error generating token when encoding payload to string: %w", err)
	}
	payload64 := base64.StdEncoding.EncodeToString(payloadstr)

	// Now add the encoded string.
	message := header64 + "." + payload64

	// We have the unsigned message ready. This is simply concat of header and payload
	unsignedStr := header + string(payloadstr)

	// we write this to the SHA256 to hash it. We can use this to generate the signature now
	h.Write([]byte(unsignedStr))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	//Finally we have the token
	tokenStr := message + "." + signature
	return tokenStr, nil
}

// This helps in validating the token
func ValidateToken(token string, secret string) error {
	// JWT has 3 parts separated by '.'
	splitToken := strings.Split(token, ".")
	// if length is not 3, we know that the token is corrupt
	if len(splitToken) != 3 {
		return errors.New(TOKEN_IS_CORRUPT)
	}

	// decode the header and payload back to strings
	header, err := base64.StdEncoding.DecodeString(splitToken[0])
	if err != nil {
		return err
	}
	payload, err := base64.StdEncoding.DecodeString(splitToken[1])
	if err != nil {
		return err
	}

	//again create the signature
	unsignedStr := string(header) + string(payload)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsignedStr))

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// if both the signature dont match, this means token is wrong
	if signature != splitToken[2] {
		return errors.New(INVALID_TOKEN)
	}

	//Unmarshal payload into ClaimsMap struct
	var payloadMap ClaimsMap
	json.Unmarshal(payload, &payloadMap)

	//Check if token is expired
	if payloadMap.Exp < fmt.Sprint(time.Now().Unix()) {
		return errors.New(TOKEN_EXPIRED)
	}

	// This means the token matches
	return nil
}
