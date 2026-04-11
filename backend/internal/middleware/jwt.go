package middleware

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//JWTSecret is the key used to sign tokens.
// IN production this comes from an environment variable - never hardcode.

var JWTSecret = []byte("pulseboard-dev-secret-change-me")

// GenerateToken creates a sogned JWT for a given user ID
func GenerateToken(userID int) (string, error) {
	//claims are the daya stored inside the token
	claims := jwt.MapClaims{
		"user_id": userID,                                //who this token belong to
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // will expire in 24 hours
		"iat":     time.Now().Unix(),                     // issued at-current time
	}

	// create a new token object with HS256 signing method and our claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ValidateToken parses a JWT string and retirmd claims if valid
func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	//jwt.Parse reads the token string and verifies it
	// the second argument is a function that returns the secret key
	// jwt calls this function to get the key for signature verification

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		//check that the signing method is what we expect (HMAC)
		// this prevents an attack where someone changes the algorith to "none"
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", t.Header["alg"])
		}
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}
	//extract the claims from the parsed token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
