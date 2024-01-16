package main

import (
	"github.com/dalefengs/chat-api-proxy/core"
	"github.com/dalefengs/chat-api-proxy/global"
	"github.com/dalefengs/chat-api-proxy/initialize"
)

func main() {
	global.Viper = core.Viper()          // 初始化Viper
	global.Log = core.Zap()              // 初始化zap日志库
	global.SugarLog = global.Log.Sugar() // zap SugarLog
	initialize.InitUserHomeCacheDir()    // 初始化文件缓存目录
	core.RunServer()
}
