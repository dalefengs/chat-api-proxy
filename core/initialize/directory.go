package initialize

import (
	"fmt"
	"github.com/dalefengs/chat-api-proxy/global"
	"os"
)

// InitUserHomeCacheDir 初始化家目录缓存目录
func InitUserHomeCacheDir() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to get user home directory")
		return
	}
	cacheDir := userHomeDir + "/.cache/poolToken"
	// 检查目录是否存在
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
			fmt.Println("初始化缓存目录失败:", err, cacheDir)
			return
		}
	} else {
		fmt.Println("缓存目录已存在:", cacheDir)
	}
	fmt.Println("初始化 Token 缓存目录成功:", cacheDir)
	global.UserHomeDir = userHomeDir
	global.UserHomeCacheDir = cacheDir
}
