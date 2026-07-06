package utils

import (
	"errors"
	"unicode/utf8"
	"strings"
)


// alphabet is the exact character set layout for standard Base62 encoding
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// EncodeBase62 takes a large numeric sequence ID and converts it into a short text token string
func EncodeBase62(id uint64) string{
	if id == 0 {
		return string(alphabet[0])
	}

	var sb strings.Builder
	base := uint64(len(alphabet)) //62

	// Repeatedly divide by 62 and collect the mathematical remainders
	for id > 0 {
		remainder := id % base
		sb.WriteByte(alphabet[remainder])
		id = id / base
	}

	// Reverse the compiled bytes to read left-to-right correctly
	result := sb.String()
	return reverseString(result)
}

// DecodeBase62 takes a short text token string and decodes it back into our database ID number
func DecodeBase62(token string) (uint64, error) {
	var id uint64
	base := uint64(len(alphabet)) // 62

	for i := 0; i < len(token); i++ {
		pos := strings.IndexByte(alphabet, token[i])
		if pos == -1 {
			return 0, errors.New("detected illegal non-alphanumeric character within the short token")
		}
		id = id*base + uint64(pos)
	}

	return id, nil
}

// reverseString flips a text string layout around backward
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}