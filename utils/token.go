package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

// GetAuthToken header 获取 token
func GetAuthToken(c *gin.Context, prefix string) (token string, err error) {
	if prefix == "" {
		prefix = "Bearer"
	}
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, prefix+" ") {
		err = errors.New("the token is invalid")
		return
	}
	prefixLen := len(prefix) + 1
	token = strings.TrimSpace(tokenString[prefixLen:])
	if token == "" {
		err = errors.New("the token is empty")
		return
	}
	return
}
