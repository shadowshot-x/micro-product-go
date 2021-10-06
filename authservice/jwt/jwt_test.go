package jwt

import (
	"fmt"
	"testing"
	"time"
)

func TestTokenValidation(t *testing.T) {

	secret := GetSecret()
	longExpiryClaims := ClaimsMap{
		Aud: "frontend.knowsearch.ml",
		Iss: "knowsearch.ml",
		Exp: fmt.Sprint(time.Now().Add(time.Minute * 60).Unix()),
	}
	longExpiryToken, err := GenerateToken("HS256", longExpiryClaims, secret)
	if err != nil {
		t.Error("Token generation failed")
	}
	//Token with long expiry date must not be expired
	if TOKEN_EXPIRED == fmt.Sprint(ValidateToken(longExpiryToken, secret)) {
		t.Error("Token must not be expired")
	}

	//Corrupt token i.e without 3 sections must throw 'Token is Corrupt' on validation
	corruptTokenString := "randomcorrupttokenstring"
	if TOKEN_IS_CORRUPT != fmt.Sprint(ValidateToken(corruptTokenString, secret)) {
		t.Error("Should throw 'Token is corrupt' for corrupt tokens")
	}

	//Invalid token i.e signature mismatched token must throw 'Invalid Token' on validation
	invalidTokenString := longExpiryToken + "randomsignaturesuffix"
	if INVALID_TOKEN != fmt.Sprint(ValidateToken(invalidTokenString, secret)) {
		t.Error("Should throw 'Token is invalid' for invalid tokens")
	}

	shortExpiryClaims := ClaimsMap{
		Aud: "frontend.knowsearch.ml",
		Iss: "knowsearch.ml",
		Exp: fmt.Sprint(time.Now().Unix()),
	}
	shortExpiryToken, err := GenerateToken("HS256", shortExpiryClaims, secret)
	if err != nil {
		t.Error("Token generation failed")
	}
	//Sleep for 5 seconds to ensure token is expired
	time.Sleep(5 * time.Second)

	//Expired token must throw 'Token Expired' on validation
	if TOKEN_EXPIRED != fmt.Sprint(ValidateToken(shortExpiryToken, secret)) {
		t.Error("Failed to detect expired token")
	}

}
