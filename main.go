package main

import (
	"github.com/dalefengs/chat-api-proxy/core"
	"github.com/dalefengs/chat-api-proxy/global"
)

var version string

func main() {
	global.Viper = core.Viper()          // 初始化Viper
	global.Log = core.Zap()              // 初始化zap日志库
	global.SugarLog = global.Log.Sugar() // zap SugarLog
	core.RunServer(version)
}
