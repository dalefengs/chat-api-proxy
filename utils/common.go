package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenHexStr Generate a random hexadecimal string of the specified length.
func GenHexStr(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
