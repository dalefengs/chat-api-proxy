package initialize

import (
	"github.com/dalefengs/chat-api-proxy/global"
	"github.com/dalefengs/chat-api-proxy/model/common/response"
	"github.com/dalefengs/chat-api-proxy/router"
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// 初始化总路由

func Routers() *gin.Engine {
	Router := gin.Default()
	copilotRouter := router.GroupApp.Copilot
	chatgptRouter := router.GroupApp.ChatGPT
	geminiRouter := router.GroupApp.GeMini
	// Router.Use(middleware.LoadTls())  // 如果需要使用https 请打开此中间件 然后前往 core/server.go 将启动模式 更变为 Router.RunTLS("端口","你的cre/pem文件","你的key文件")
	// 跨域，如需跨域可以打开下面的注释
	// Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	// Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	//global.Log.Info("use middleware cors")
	Router.GET(global.Config.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.Log.Info("register swagger handler")

	Router.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		e := err.(error)
		response.FailWithMessage(e.Error(), c)
		return
	}))
	// 方便统一添加路由组前缀 多服务器上线使用

	PublicGroup := Router.Group(global.Config.System.RouterPrefix)
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
	{
		copilotRouter.InitCopilotRouter(PublicGroup)
		copilotRouter.InitCoCopilotRouter(PublicGroup)
		chatgptRouter.InitChatGPTRouter(PublicGroup)
		geminiRouter.InitGeminiRouter(PublicGroup)
	}
	global.Log.Info("router register success")
	return Router
}
