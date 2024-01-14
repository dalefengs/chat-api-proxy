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

func SetEventStreamHeaders(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
}
