package copilot

import (
	"github.com/dalefengs/chat-api-proxy/api"
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (s *Router) InitCopilotRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("copilot")
	copilotApi := api.ApiGroupApp.CopilotApi
	{
		baseRouter.GET("/copilot_internal/v2/token", copilotApi.TokenHandler) // 官方获取 token
		baseRouter.POST("/v1/chat/completions", copilotApi.CompletionsHandler)
	}
	return baseRouter
}

func (s *Router) InitCoCopilotRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("cocopilot")
	copilotApi := api.ApiGroupApp.CopilotApi
	{
		baseRouter.GET("/copilot_internal/v2/token", copilotApi.CoTokenHandler)
	}
	return baseRouter
}
