package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

// GetAuthToken header 获取 token
func GetAuthToken(c *gin.Context) (token string, err error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		err = errors.New("the token is invalid")
		return
	}
	token = tokenString[7:]
	return
}
