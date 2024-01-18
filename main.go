package main

import (
	"fmt"
	"github.com/dalefengs/chat-api-proxy/core"
	"github.com/dalefengs/chat-api-proxy/global"
)

var version string

func main() {
	global.Viper = core.Viper()          // 初始化Viper
	global.Log = core.Zap()              // 初始化zap日志库
	global.SugarLog = global.Log.Sugar() // zap SugarLog
	core.InitEnvironmentVar()
	fmt.Println("version", version)
	core.RunServer(version)
}
