package copilot

import (
	"github.com/dalefeng/chat-api-reverse/api"
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (s *Router) InitCopilotRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("copilot")
	copilotApi := api.ApiGroupApp.CopilotApi
	{
		baseRouter.GET("copilot_internal/v2/token", copilotApi.CoToken) // 获取 token
		baseRouter.POST("v1/chat/completions", copilotApi.CoCompletions)
	}
	return baseRouter
}

func (s *Router) InitCoCopilotRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("cocopilot")
	copilotApi := api.ApiGroupApp.CopilotApi
	{
		baseRouter.GET("copilot_internal/v2/token", copilotApi.CoToken) // 获取 token
		baseRouter.POST("v1/chat/completions", copilotApi.CoCompletions)
	}
	return baseRouter
}
