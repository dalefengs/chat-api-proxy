package main

import (
	"github.com/dalefeng/chat-api-reverse/core"
	"github.com/dalefeng/chat-api-reverse/global"
)

func main() {
	global.Viper = core.Viper()          // 初始化Viper
	global.Log = core.Zap()              // 初始化zap日志库
	global.SugarLog = global.Log.Sugar() // zap SugarLog

	core.RunServer()
}
