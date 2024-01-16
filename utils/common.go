package utils

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/dalefengs/chat-api-proxy/global"
)

// GenHexStr Generate a random hexadecimal string of the specified length.
func GenHexStr(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func Int2Byte(num int) []byte {
	return []byte{byte(num >> 24), byte(num >> 16), byte(num >> 8), byte(num)}
}

func Byte2Int(b []byte) int {
	return int(b[3]) | int(b[2])<<8 | int(b[1])<<16 | int(b[0])<<24
}

// GetTokenCacheFilePath Get the path of the token cache file.
func GetTokenCacheFilePath(token string) string {
	return global.UserHomeCacheDir + "/" + token
}
