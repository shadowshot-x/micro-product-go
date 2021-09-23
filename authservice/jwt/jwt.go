package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// Function for generating the tokens.
func GenerateToken(header string, payload map[string]string, secret string) (string, error) {
	// create a new hash of type sha256. We pass the secret key to it
	// sha256 is a symmetric cryptographic algorithm
	h := hmac.New(sha256.New, []byte(secret))

	// We base encode the header which is a normal string
	header64 := base64.StdEncoding.EncodeToString([]byte(header))
	// We then Marshal the payload which is a map. This converts it to a string of JSON.
	// Now we base encode this string
	payloadstr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error generating Token")
		return string(payloadstr), err
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
func ValidateToken(token string, secret string) (bool, error) {
	// JWT has 3 parts separated by '.'
	splitToken := strings.Split(token, ".")
	// if length is not 3, we know that the token is corrupt
	if len(splitToken) != 3 {
		return false, nil
	}

	// decode the header and payload back to strings
	header, err := base64.StdEncoding.DecodeString(splitToken[0])
	if err != nil {
		return false, err
	}
	payload, err := base64.StdEncoding.DecodeString(splitToken[1])
	if err != nil {
		return false, err
	}

	//again create the signature
	unsignedStr := string(header) + string(payload)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsignedStr))

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Println(signature)

	// if both the signature dont match, this means token is wrong
	if signature != splitToken[2] {
		return false, nil
	}
	// This means the token matches
	return true, nil
}
