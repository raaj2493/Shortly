package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)


// GenerateToken creates a signed digital key for a logged-in user that lasts 24 hours
func GenerateToken(userID string, secret string) (string, error) {
	// 1. Create the payload data (Claims)
	claims := jwt.MapClaims{
		"sub": userID,                                 // "Subject" - who owns this token
		"exp": time.Now().Add(24 * time.Hour).Unix(), // Expiration time (24 hours from now)
		"iat": time.Now().Unix(),                      // Issued At time
	}

	// 2. Choose the signing algorithm (HMAC-SHA256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. Sign it with our secret key
	return token.SignedString([]byte(secret))
}