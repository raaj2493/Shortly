package utils


import "golang.org/x/crypto/bcrypt"

func HashPass( password string )(string , error ){
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPasswordHash compares a raw login password with the saved database hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // returns true if they match, false if they don't!
}