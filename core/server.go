package core

import (
	"fmt"
	"github.com/dalefengs/chat-api-proxy/core/initialize"
	"github.com/dalefengs/chat-api-proxy/global"
	"runtime"
	"time"

	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunServer(version string) {
	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	port := fmt.Sprintf(":%d", global.Config.System.Port)
	s := initServer(port, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.Log.Info("server run success on ", zap.String("address", port))

	fmt.Printf(`System info:
================================================================================
=
=  OS:          %s
=  Golang:      %s
=  Version:     %s
=  Bind:        http://127.0.0.1:%d/
=
================================================================================`+"\n", runtime.GOOS, runtime.Version(), version, global.Config.System.Port)
	global.Log.Error(s.ListenAndServe().Error())
}
