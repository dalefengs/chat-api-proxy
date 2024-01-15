package core

import (
	"fmt"
	"github.com/dalefengs/chat-api-proxy/global"
	"github.com/dalefengs/chat-api-proxy/initialize"
	"time"

	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.Config.System.Port)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.Log.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`System info:
================================================================================
=
=  OS:          Linux
=  Prefix:      Linux
=  Bind:        http://0.0.0.0:%s/
=
================================================================================`+"\n", address)
	global.Log.Error(s.ListenAndServe().Error())
}
